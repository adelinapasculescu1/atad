package importsvc

import (
    "encoding/csv"
    "fmt"
    "io"
    "os"
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
        // TODO: implement OFX support
        return fmt.Errorf("OFX format not implemented yet")
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
