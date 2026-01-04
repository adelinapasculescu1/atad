package models

type Budget struct {
    ID         int64
    CategoryID int64
    Amount     float64
    Period     string 
}
