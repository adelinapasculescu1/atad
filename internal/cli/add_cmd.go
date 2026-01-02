package cli

import (
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/adelinapasculescu1/atad/internal/services/transactionsvc"
	"github.com/adelinapasculescu1/atad/internal/models"
    "github.com/spf13/cobra"
)

func NewAddCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "add",
        Short: "Add a manual income or expense transaction",
        RunE: func(cmd *cobra.Command, args []string) error {
            amountStr, _ := cmd.Flags().GetString("amount")
            typeStr, _ := cmd.Flags().GetString("type")
            desc, _ := cmd.Flags().GetString("description")
            dateStr, _ := cmd.Flags().GetString("date")
            categoryName, _ := cmd.Flags().GetString("category")
            
            if amountStr == "" {
                return fmt.Errorf("--amount is required")
            }
            if typeStr == "" {
                return fmt.Errorf("--type is required (income|expense)")
            }
            if desc == "" {
                return fmt.Errorf("--description is required")
            }

            amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
            if err != nil {
                return fmt.Errorf("invalid amount %q: %w", amountStr, err)
            }

            var txType models.TransactionType
            switch strings.ToLower(strings.TrimSpace(typeStr)) {
            case "income":
                txType = models.TransactionTypeIncome
            case "expense":
                txType = models.TransactionTypeExpense
            default:
                return fmt.Errorf("invalid type %q, expected income or expense", typeStr)
            }

            var date time.Time
            if dateStr == "" {
                date = time.Now()
            } else {
                // Format YYYY-MM-DD
                d, err := time.Parse("2006-01-02", strings.TrimSpace(dateStr))
                if err != nil {
                    return fmt.Errorf("invalid date %q, expected YYYY-MM-DD: %w", dateStr, err)
                }
                date = d
            }

            var categoryID *int64
            if strings.TrimSpace(categoryName) != "" {
                cat, err := deps.CategoryRepo.GetByName(strings.TrimSpace(categoryName))
                if err != nil {
                    return err
                }
                if cat == nil {
                    return fmt.Errorf("category not found: %s (create it with `atad category add --name %q`)", categoryName, categoryName)
                }
                categoryID = &cat.ID
            }


            input := transactionsvc.AddManualInput{
                Date:        date,
                Description: desc,
                Amount:      amount,
                Type:        txType,
                CategoryID:  categoryID,

            }

            if err := deps.TransactionSvc.AddManual(input); err != nil {
                return err
            }

            fmt.Println("Transaction added successfully.")
            return nil
        },
    }

    cmd.Flags().String("amount", "", "Transaction amount (e.g. 50.0)")
    cmd.Flags().String("type", "", "Transaction type: income or expense")
    cmd.Flags().String("description", "", "Transaction description")
    cmd.Flags().String("date", "", "Transaction date (YYYY-MM-DD). Defaults to today.")
    cmd.Flags().String("category", "", "Category name (optional)")


    return cmd
}
