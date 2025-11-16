package main

import (
    "log"

    "github.com/adelinapasculescu1/atad/internal/cli"
    "github.com/adelinapasculescu1/atad/internal/db"
    "github.com/adelinapasculescu1/atad/internal/repository"
    "github.com/adelinapasculescu1/atad/internal/services/importsvc"
)


func main() {
    // Pentru început hardcodăm calea DB
    dbPath := "atad.db"

    database, err := db.Open(dbPath)
    if err != nil {
        log.Fatalf("error opening database: %v", err)
    }
    defer database.Close()

    if err := db.Migrate(database); err != nil {
        log.Fatalf("error running migrations: %v", err)
    }

    // Repo & services
    txRepo := repository.NewSQLiteTransactionRepository(database)
    importService := importsvc.NewService(txRepo)

    deps := cli.Deps{
        DB:        database,
        ImportSvc: importService,
    }

    rootCmd := cli.NewRootCommand(deps)
    if err := rootCmd.Execute(); err != nil {
        log.Fatalf("command error: %v", err)
    }
}

