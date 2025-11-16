package cli

import (
    "fmt"

    "github.com/spf13/cobra"
)

func NewImportCommand(deps Deps) *cobra.Command {
    cmd := &cobra.Command{
        Use:   "import",
        Short: "Import transactions from a CSV/OFX file",
        RunE: func(cmd *cobra.Command, args []string) error {
            file, _ := cmd.Flags().GetString("file")
            format, _ := cmd.Flags().GetString("format")

            if file == "" {
                return fmt.Errorf("--file is required")
            }
            if format == "" {
                format = "csv"
            }

            return deps.ImportSvc.Import(file, format)
        },
    }

    cmd.Flags().StringP("file", "f", "", "Path to CSV/OFX file")
    cmd.Flags().StringP("format", "", "csv", "File format: csv or ofx")

    return cmd
}
