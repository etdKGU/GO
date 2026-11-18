package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"html/template"
	"myserver/internal/repository"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

var Users *repository.UserRepository
var Sessions *repository.SessionRepository

type registerPage struct {
	Error string
	Email string
}

func renderRegister(w http.ResponseWriter, status int, data registerPage) {
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/register.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	tmpl.ExecuteTemplate(w, "layout", data)
}

func RegisterForm(w http.ResponseWriter, r *http.Request) {
	renderRegister(w, http.StatusOK, registerPage{})
}

func Register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	if !strings.Contains(email, "@") {
		renderRegister(w, http.StatusBadRequest, registerPage{Error: "Некорректный email", Email: email})
		return
	}
	if len(password) < 6 || len(password) > 72 {
		renderRegister(w, http.StatusBadRequest, registerPage{Error: "Пароль должен быть от 6 до 72 символов", Email: email})
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	err = Users.Create(email, string(hash))
	if err != nil {
		renderRegister(w, http.StatusBadRequest, registerPage{Error: "Такой email уже зарегистрирован", Email: email})
		return
	}

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func newSessionID() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func renderLogin(w http.ResponseWriter, status int, data registerPage) {
	tmpl, err := template.ParseFiles("web/templates/layout.html", "web/templates/login.html")
	if err != nil {
		http.Error(w, "Ошибка шаблона", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(status)
	tmpl.ExecuteTemplate(w, "layout", data)
}

func LoginForm(w http.ResponseWriter, r *http.Request) {
	renderLogin(w, http.StatusOK, registerPage{})
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "Ошибка формы", http.StatusBadRequest)
		return
	}

	email := strings.TrimSpace(r.FormValue("email"))
	password := r.FormValue("password")

	user, err := Users.GetByEmail(email)
	if err != nil || bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)) != nil {
		renderLogin(w, http.StatusUnauthorized, registerPage{Error: "Неверный email или пароль", Email: email})
		return
	}

	sessionID, err := newSessionID()
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	expires := time.Now().Add(24 * time.Hour)
	err = Sessions.Create(sessionID, user.ID, expires)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		Expires:  expires,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/expenses", http.StatusSeeOther)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err == nil {
		Sessions.Delete(cookie.Value)
	}

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
