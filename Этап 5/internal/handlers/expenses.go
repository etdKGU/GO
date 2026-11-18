package handlers

import (
	"html/template"
	"myserver/internal/models"
	"myserver/internal/repository"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var Repo *repository.ExpenseRepository

func List(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	categoryID, _ := strconv.Atoi(q.Get("category"))

	filter := models.ExpenseFilter{
		CategoryID: categoryID,
		From:       q.Get("from"),
		To:         q.Get("to"),
		UserID:     UserID(r),
	}

	items, err := Repo.GetAll(filter)
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	categories, err := Repo.GetCategories()
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/list.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := struct {
		Expenses   []models.Expense
		Categories []models.Category
		Filter     models.ExpenseFilter
	}{
		Expenses:   items,
		Categories: categories,
		Filter:     filter,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func NewForm(w http.ResponseWriter, r *http.Request) {
	categories, err := Repo.GetCategories()
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/form.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := struct {
		Categories []models.Category
	}{
		Categories: categories,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func Create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	amountStr := r.FormValue("amount")
	description := strings.TrimSpace(r.FormValue("description"))
	date := r.FormValue("date")

	amount, err := strconv.ParseFloat(amountStr, 64)
	if err != nil || amount <= 0 {
		http.Error(w, "Неверная сумма", http.StatusBadRequest)
		return
	}

	if description == "" {
		http.Error(w, "Поле должно быть заполненным", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		http.Error(w, "Некорректная дата", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Неверная категория", http.StatusBadRequest)
		return
	}

	e := models.Expense{
		Amount:      amount,
		Description: description,
		Date:        date,
		CategoryID:  categoryID,
		UserID:      UserID(r),
	}

	err = Repo.Create(e)
	if err != nil {
		http.Error(w, "Не удалось сохранить трату", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func Delete(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Неверный номер траты", http.StatusBadRequest)
		return
	}

	err = Repo.Delete(id, UserID(r))
	if err != nil {
		http.Error(w, "Не удалось удалить трату", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func EditForm(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil {
		http.Error(w, "Неверный номер траты", http.StatusBadRequest)
		return
	}

	e, err := Repo.GetByID(id, UserID(r))
	if err != nil {
		http.Error(w, "Трата не найдена", http.StatusNotFound)
		return
	}

	categories, err := Repo.GetCategories()
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/edit.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}

	data := struct {
		models.Expense
		Categories []models.Category
	}{
		Expense:    e,
		Categories: categories,
	}
	tmpl.ExecuteTemplate(w, "layout", data)
}

func Update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Неверный номер траты", http.StatusBadRequest)
		return
	}

	amountStr := r.FormValue("amount")
	description := strings.TrimSpace(r.FormValue("description"))
	date := r.FormValue("date")
	amount, err := strconv.ParseFloat(amountStr, 64)

	if err != nil || amount <= 0 {
		http.Error(w, "Неверная сумма", http.StatusBadRequest)
		return
	}

	if description == "" {
		http.Error(w, "Поле должно быть заполненным", http.StatusBadRequest)
		return
	}

	_, err = time.Parse("2006-01-02", date)
	if err != nil {
		http.Error(w, "Некорректная дата", http.StatusBadRequest)
		return
	}

	categoryID, err := strconv.Atoi(r.FormValue("category_id"))
	if err != nil {
		http.Error(w, "Неверная категория", http.StatusBadRequest)
		return
	}

	e := models.Expense{
		ID:          id,
		Amount:      amount,
		Description: description,
		Date:        date,
		CategoryID:  categoryID,
		UserID:      UserID(r),
	}

	err = Repo.Update(e)
	if err != nil {
		http.Error(w, "Не удалось сохранить изменения", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}
