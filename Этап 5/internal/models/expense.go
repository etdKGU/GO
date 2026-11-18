package models

type Expense struct {
	ID           int
	Amount       float64
	Description  string
	Date         string
	CategoryID   int
	CategoryName string
	UserID       int
}

type ExpenseFilter struct {
	CategoryID int
	From       string
	To         string
	UserID     int
}
