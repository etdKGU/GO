package repository

import (
	"database/sql"
	"myserver/internal/models"
	"strings"
)

type ExpenseRepository struct {
	db *sql.DB
}

func NewExpenseRepository(db *sql.DB) *ExpenseRepository {
	return &ExpenseRepository{db: db}
}

func (r *ExpenseRepository) Create(e models.Expense) error {
	_, err := r.db.Exec(
		"INSERT INTO expenses (amount, description, date, category_id) VALUES (?, ?, ?, ?)",
		e.Amount, e.Description, e.Date, e.CategoryID,
	)
	return err
}

func (r *ExpenseRepository) Update(e models.Expense) error {
	_, err := r.db.Exec(
		"UPDATE expenses SET amount = ?, description = ?, date = ?, category_id = ? WHERE id = ?",
		e.Amount, e.Description, e.Date, e.CategoryID, e.ID,
	)
	return err
}

func (r *ExpenseRepository) GetByID(id int) (models.Expense, error) {
	var e models.Expense
	err := r.db.QueryRow(
		"SELECT id, amount, description, date, category_id FROM expenses WHERE id = ?",
		id,
	).Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID)
	return e, err
}

func (r *ExpenseRepository) GetAll(f models.ExpenseFilter) ([]models.Expense, error) {
	query := `
        SELECT e.id, e.amount, e.description, e.date, e.category_id, c.name
        FROM expenses e
        JOIN categories c ON c.id = e.category_id`

	var conditions []string
	var args []any

	if f.CategoryID > 0 {
		conditions = append(conditions, "e.category_id = ?")
		args = append(args, f.CategoryID)
	}
	if f.From != "" {
		conditions = append(conditions, "e.date >= ?")
		args = append(args, f.From)
	}
	if f.To != "" {
		conditions = append(conditions, "e.date <= ?")
		args = append(args, f.To)
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	query += " ORDER BY e.id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Expense
	for rows.Next() {
		var e models.Expense
		err := rows.Scan(&e.ID, &e.Amount, &e.Description, &e.Date, &e.CategoryID, &e.CategoryName)
		if err != nil {
			return nil, err
		}
		result = append(result, e)
	}
	return result, rows.Err()
}

func (r *ExpenseRepository) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM expenses WHERE id = ?", id)
	return err
}

func (r *ExpenseRepository) GetCategories() ([]models.Category, error) {
	rows, err := r.db.Query("SELECT id, name FROM categories ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []models.Category
	for rows.Next() {
		var c models.Category
		err := rows.Scan(&c.ID, &c.Name)
		if err != nil {
			return nil, err
		}
		result = append(result, c)
	}
	return result, rows.Err()
}
