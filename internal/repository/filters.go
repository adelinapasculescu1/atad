package repository

import "time"

type TransactionFilter struct {
    From *time.Time
    To   *time.Time

    Type *string 

    CategoryName *string 

    Query *string 

    MinAmount *float64
    MaxAmount *float64

    Limit  int
    Offset int
}
