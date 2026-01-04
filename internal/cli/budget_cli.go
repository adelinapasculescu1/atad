package cli

import (
    "fmt"
    "strconv"
    "strings"
    "time"

    "github.com/spf13/cobra"
)

func NewBudgetCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "budget",
        Short: "Manage budgets",
    }

    cmd.AddCommand(NewBudgetSetCommand(deps))
    cmd.AddCommand(NewBudgetStatusCommand(deps))
    return cmd
}

func NewBudgetSetCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "set",
        Short: "Set a budget for a category",
        RunE: func(cmd *cobra.Command, args []string) error {
            category, _ := cmd.Flags().GetString("category")
            amountStr, _ := cmd.Flags().GetString("amount")
            period, _ := cmd.Flags().GetString("period")

            if strings.TrimSpace(category) == "" {
                return fmt.Errorf("--category is required")
            }
            if strings.TrimSpace(amountStr) == "" {
                return fmt.Errorf("--amount is required")
            }

            amount, err := strconv.ParseFloat(strings.TrimSpace(amountStr), 64)
            if err != nil {
                return fmt.Errorf("invalid amount: %w", err)
            }

            if err := deps.BudgetSvc.SetBudget(category, amount, period); err != nil {
                return err
            }

            fmt.Println("Budget set successfully.")
            return nil
        },
    }

    cmd.Flags().String("category", "", "Category name")
    cmd.Flags().String("amount", "", "Budget amount (e.g. 500)")
    cmd.Flags().String("period", "monthly", "Budget period (only monthly supported)")
    return cmd
}

func NewBudgetStatusCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "status",
        Short: "Show budget status for a given month",
        RunE: func(cmd *cobra.Command, args []string) error {
            monthStr, _ := cmd.Flags().GetString("month")
            if strings.TrimSpace(monthStr) == "" {
                // default: luna curenta
                now := time.Now()
                monthStr = fmt.Sprintf("%04d-%02d", now.Year(), int(now.Month()))
            }

            t, err := time.Parse("2006-01", monthStr)
            if err != nil {
                return fmt.Errorf("invalid --month, expected YYYY-MM: %w", err)
            }

            statuses, err := deps.BudgetSvc.GetMonthlyStatus(t.Year(), t.Month())
            if err != nil {
                return err
            }

            if len(statuses) == 0 {
                fmt.Println("No budgets found.")
                return nil
            }

            for _, s := range statuses {
                fmt.Printf(
                    "%s\tBudget: %.2f\tSpent: %.2f\tRemaining: %.2f\t%s\n",
                    s.CategoryName, s.BudgetAmount, s.SpentAmount, s.Remaining, s.Status,
                )
            }

            return nil
        },
    }

    cmd.Flags().String("month", "", "Month in YYYY-MM (default: current month)")
    return cmd
}
