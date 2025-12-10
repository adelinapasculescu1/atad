package importsvc

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
    "regexp"
    "strconv"
    "strings"
    "time"

    "github.com/adelinapasculescu1/atad/internal/models"
    "github.com/adelinapasculescu1/atad/internal/repository"
)

type Service struct {
    txRepo repository.TransactionRepository
}

// Constructor
func NewService(txRepo repository.TransactionRepository) *Service {
    return &Service{txRepo: txRepo}
}

// csv or ofx, ofx to be added
func (s *Service) Import(filePath string, format string) error {
    switch strings.ToLower(format) {
    case "csv":
        return s.importCSV(filePath)
    case "ofx":
        return s.importOFX(filePath)
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}

func (s *Service) importCSV(filePath string) error {
    f, err := os.Open(filePath)
    if err != nil {
        return err
    }
    defer f.Close()

    reader := csv.NewReader(f)
    reader.TrimLeadingSpace = true

    header, err := reader.Read()
    if err != nil {
        return fmt.Errorf("reading header: %w", err)
    }

    colIndex := map[string]int{}
    for i, col := range header {
        colLower := strings.ToLower(strings.TrimSpace(col))
        colIndex[colLower] = i
    }

    required := []string{"date", "description", "amount", "type"}
    for _, col := range required {
        if _, ok := colIndex[col]; !ok {
            return fmt.Errorf("missing required column %q in CSV", col)
        }
    }

    importedCount := 0

    for {
        record, err := reader.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            return fmt.Errorf("reading record: %w", err)
        }

        dateStr := record[colIndex["date"]]
        desc := record[colIndex["description"]]
        amountStr := record[colIndex["amount"]]
        typeStr := strings.ToLower(strings.TrimSpace(record[colIndex["type"]]))

        //format YYYY-MM-DD
        date, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
        if err != nil {
            return fmt.Errorf("invalid date %q: %w", dateStr, err)
        }

        amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
        if err != nil {
            return fmt.Errorf("invalid amount %q: %w", amountStr, err)
        }

        var txType models.TransactionType
        switch typeStr {
        case "income":
            txType = models.TransactionTypeIncome
        case "expense":
            txType = models.TransactionTypeExpense
        default:
            if amount >= 0 {
                txType = models.TransactionTypeIncome
            } else {
                txType = models.TransactionTypeExpense
            }
        }

        tx := &models.Transaction{
            Date:        date,
            Description: desc,
            Amount:      amount,
            Type:        txType,
            CategoryID:  nil, 
            Source:      fmt.Sprintf("import:%s", filePath),
        }

        if err := s.txRepo.Create(tx); err != nil {
            return fmt.Errorf("saving transaction: %w", err)
        }

        importedCount++
    }

    fmt.Printf("Imported %d transactions from %s\n", importedCount, filePath)
    return nil
}

func (s *Service) importOFX(filePath string) error {
    data, err := os.ReadFile(filePath)
    if err != nil {
        return fmt.Errorf("opening OFX file: %w", err)
    }

    content := string(data)

    // Căutăm blocurile <STMTTRN> ... </STMTTRN>
    re := regexp.MustCompile(`(?s)<STMTTRN>(.*?)</STMTTRN>`)
    matches := re.FindAllStringSubmatch(content, -1)
    if len(matches) == 0 {
        return fmt.Errorf("no <STMTTRN> blocks found in OFX file")
    }

    importedCount := 0

    for _, m := range matches {
        block := m[1]

        rawDate := getOFXTagValue(block, "DTPOSTED")
        rawAmount := getOFXTagValue(block, "TRNAMT")
        rawType := strings.ToUpper(strings.TrimSpace(getOFXTagValue(block, "TRNTYPE")))
        desc := getOFXTagValue(block, "NAME")
        if desc == "" {
            desc = getOFXTagValue(block, "MEMO")
        }

        if rawDate == "" || rawAmount == "" {
            // Sărim peste tranzacții incomplete
            continue
        }

        // DTPOSTED de forma 20250118120000 → luăm doar YYYYMMDD
        if len(rawDate) < 8 {
            return fmt.Errorf("invalid DTPOSTED value: %q", rawDate)
        }
        datePart := rawDate[:8]

        date, err := time.Parse("20060102", datePart)
        if err != nil {
            return fmt.Errorf("invalid date %q: %w", rawDate, err)
        }

        amount, err := strconv.ParseFloat(strings.TrimSpace(rawAmount), 64)
        if err != nil {
            return fmt.Errorf("invalid amount %q: %w", rawAmount, err)
        }

        // Determinăm tipul tranzacției
        txType := models.TransactionTypeExpense
        if amount > 0 {
            txType = models.TransactionTypeIncome
        }
        // Dacă există TRNTYPE, îl folosim ca hint
        if rawType == "CREDIT" {
            txType = models.TransactionTypeIncome
        } else if rawType == "DEBIT" {
            txType = models.TransactionTypeExpense
        }

        tx := &models.Transaction{
            Date:        date,
            Description: desc,
            Amount:      amount,
            Type:        txType,
            CategoryID:  nil,
            Source:      fmt.Sprintf("import:%s", filePath),
        }

        if err := s.txRepo.Create(tx); err != nil {
            return fmt.Errorf("saving transaction: %w", err)
        }

        importedCount++
    }

    fmt.Printf("Imported %d OFX transactions from %s\n", importedCount, filePath)
    return nil
}

// getOFXTagValue caută <TAG>valoare</TAG> sau <TAG>valoare<următor-tag>
// și returnează "valoare".
func getOFXTagValue(block, tag string) string {
    tagStart := "<" + tag + ">"
    idx := strings.Index(block, tagStart)
    if idx == -1 {
        return ""
    }

    start := idx + len(tagStart)

    // Căutăm următorul '<' după valoare
    end := strings.Index(block[start:], "<")
    if end == -1 {
        end = len(block)
    } else {
        end = start + end
    }

    return strings.TrimSpace(block[start:end])
}


