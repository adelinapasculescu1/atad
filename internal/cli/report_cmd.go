package cli

import (
    "fmt"
    "strings"
    "time"

    "github.com/spf13/cobra"
)

func NewReportCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "report",
        Short: "Generate reports",
    }

    cmd.AddCommand(NewMonthlyReportCommand(deps))
    return cmd
}

func NewMonthlyReportCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "monthly",
        Short: "Monthly spending report",
        RunE: func(cmd *cobra.Command, args []string) error {
            monthStr, _ := cmd.Flags().GetString("month")

            if strings.TrimSpace(monthStr) == "" {
                now := time.Now()
                monthStr = fmt.Sprintf("%04d-%02d", now.Year(), int(now.Month()))
            }

            t, err := time.Parse("2006-01", monthStr)
            if err != nil {
                return fmt.Errorf("invalid month format, expected YYYY-MM")
            }

            report, err := deps.ReportSvc.Monthly(t.Year(), t.Month())
            if err != nil {
                return err
            }

            fmt.Printf("\nReport for %04d-%02d\n\n", report.Year, report.Month)
            fmt.Printf("TOTAL SPENT: %.2f\n\n", report.TotalSpent)

            if len(report.ByCategoryName) == 0 {
                fmt.Println("No expenses for this period.")
                return nil
            }

            for cat, sum := range report.ByCategoryName {
                fmt.Printf("%-15s %.2f\n", cat, sum)
            }

            return nil
        },
    }

    cmd.Flags().String("month", "", "Month in YYYY-MM (default: current month)")
    return cmd
}
