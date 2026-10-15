package handlers

import (
	"html/template"
	"myserver/internal/models"
	"net/http"
	"strconv"
)

var expenses = []models.Expense{
	{ID: 1, Amount: 250.00, Description: "Обед", Date: "2026-10-14"},
	{ID: 2, Amount: 500.00, Description: "Доставка", Date: "2026-10-15"},
	{ID: 3, Amount: 1000.00, Description: "Долг", Date: "2026-10-16"},
}

func List(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/list.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", expenses)
}

func NewForm(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/form.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", nil)
}

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	amountStr := r.FormValue("amount")
	description := r.FormValue("description")
	date := r.FormValue("date")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil {
		http.Error(w, "Неверная сумма", http.StatusBadRequest)
		return
	}

	e := models.Expense{
		ID:          len(expenses) + 1,
		Amount:      amount,
		Description: description,
		Date:        date,
	}

	expenses = append(expenses, e)

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
