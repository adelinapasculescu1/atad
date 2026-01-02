package cli

import (
    "database/sql"
    "github.com/adelinapasculescu1/atad/internal/repository"
    "github.com/adelinapasculescu1/atad/internal/services/importsvc"
    "github.com/adelinapasculescu1/atad/internal/services/transactionsvc"
    "github.com/spf13/cobra"
)


type Deps struct {
    DB          *sql.DB
    ImportSvc   *importsvc.Service
    TransactionSvc *transactionsvc.Service
    CategoryRepo repository.CategoryRepository

    //to be added more
}


func NewRootCommand(deps Deps) *cobra.Command {
    rootCmd := &cobra.Command{
        Use:   "atad",
        Short: "ATAD - command-line personal finance manager",
    }

    // subcomenzi
    rootCmd.AddCommand(
        NewImportCommand(deps),
        NewAddCommand(deps),
        NewCategoryCommand(deps),
    )

    return rootCmd
}


