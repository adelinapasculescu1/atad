package models

import "time"

type TransactionType string

const (
    TransactionTypeIncome  TransactionType = "income"
    TransactionTypeExpense TransactionType = "expense"
)

type Transaction struct {
    ID          int64
    Date        time.Time
    Description string
    Amount      float64
    Type        TransactionType
    CategoryID  *int64
    Source      string
}
