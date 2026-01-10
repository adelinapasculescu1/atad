package cli

import (
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/adelinapasculescu1/atad/internal/repository"
    "github.com/spf13/cobra"
)

func NewSearchCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "search",
        Short: "Search and filter transactions",
        RunE: func(cmd *cobra.Command, args []string) error {
            fromStr, _ := cmd.Flags().GetString("from")
            toStr, _ := cmd.Flags().GetString("to")
            typeStr, _ := cmd.Flags().GetString("type")
            categoryStr, _ := cmd.Flags().GetString("category")
            qStr, _ := cmd.Flags().GetString("q")
            minStr, _ := cmd.Flags().GetString("min")
            maxStr, _ := cmd.Flags().GetString("max")
            limit, _ := cmd.Flags().GetInt("limit")

            var filter repository.TransactionFilter
            filter.Limit = limit

            if strings.TrimSpace(fromStr) != "" {
                t, err := time.Parse("2006-01-02", strings.TrimSpace(fromStr))
                if err != nil {
                    return fmt.Errorf("invalid --from, expected YYYY-MM-DD")
                }
                filter.From = &t
            }

            if strings.TrimSpace(toStr) != "" {
                t, err := time.Parse("2006-01-02", strings.TrimSpace(toStr))
                if err != nil {
                    return fmt.Errorf("invalid --to, expected YYYY-MM-DD")
                }
                // Interpretăm --to ca inclusiv: adăugăm o zi și folosim < end
                t = t.AddDate(0, 0, 1)
                filter.To = &t
            }

            if strings.TrimSpace(typeStr) != "" {
                t := strings.ToLower(strings.TrimSpace(typeStr))
                if t != "income" && t != "expense" {
                    return fmt.Errorf("invalid --type, expected income or expense")
                }
                filter.Type = &t
            }

            if strings.TrimSpace(categoryStr) != "" {
                c := strings.TrimSpace(categoryStr)
                filter.CategoryName = &c
            }

            if strings.TrimSpace(qStr) != "" {
                q := strings.TrimSpace(qStr)
                filter.Query = &q
            }

            if strings.TrimSpace(minStr) != "" {
                v, err := strconv.ParseFloat(strings.TrimSpace(minStr), 64)
                if err != nil {
                    return fmt.Errorf("invalid --min")
                }
                filter.MinAmount = &v
            }

            if strings.TrimSpace(maxStr) != "" {
                v, err := strconv.ParseFloat(strings.TrimSpace(maxStr), 64)
                if err != nil {
                    return fmt.Errorf("invalid --max")
                }
                filter.MaxAmount = &v
            }

            results, err := deps.SearchSvc.Search(filter)
            if err != nil {
                return err
            }

            if len(results) == 0 {
                fmt.Println("No transactions found.")
                return nil
            }

            for _, tx := range results {
                fmt.Printf(
                    "%d\t%s\t%s\t%.2f\t%s\tcat=%v\t%s\n",
                    tx.ID,
                    tx.Date.Format("2006-01-02"),
                    tx.Description,
                    tx.Amount,
                    tx.Type,
                    tx.CategoryID,
                    tx.Source,
                )
            }

            return nil
        },
    }

    cmd.Flags().String("from", "", "Start date YYYY-MM-DD")
    cmd.Flags().String("to", "", "End date YYYY-MM-DD (inclusive)")
    cmd.Flags().String("type", "", "Transaction type: income|expense")
    cmd.Flags().String("category", "", "Category name")
    cmd.Flags().String("q", "", "Keyword search in description")
    cmd.Flags().String("min", "", "Minimum amount (absolute)")
    cmd.Flags().String("max", "", "Maximum amount (absolute)")
    cmd.Flags().Int("limit", 50, "Max results")

    return cmd
}
