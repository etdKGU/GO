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
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/list.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	items, err := Repo.GetAll()
	if err != nil {
		http.Error(w, "Ошибка базы данных", http.StatusInternalServerError)
		return
	}
	tmpl.ExecuteTemplate(w, "layout", items)
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

	e := models.Expense{
		Amount:      amount,
		Description: description,
		Date:        date,
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

    err = Repo.Delete(id)
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

    e, err := Repo.GetByID(id)
    if err != nil {
        http.Error(w, "Трата не найдена", http.StatusNotFound)
        return
    }

    tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/edit.html")
    if err != nil {
        http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
        return
    }

    tmpl.ExecuteTemplate(w, "layout", e)
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

    e := models.Expense{
        ID:          id,
        Amount:      amount,
        Description: description,
        Date:        date,
    }

    err = Repo.Update(e)
    if err != nil {
        http.Error(w, "Не удалось сохранить изменения", http.StatusInternalServerError)
        return
    }

    http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}