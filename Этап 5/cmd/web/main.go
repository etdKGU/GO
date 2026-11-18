package main

import (
	"database/sql"
	"fmt"
	"log"
	"myserver/internal/handlers"
	"myserver/internal/repository"
	"net/http"
	"os"

	_ "modernc.org/sqlite"
)

func sait(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, это главная страница!")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Это мой первый проект на языке GO!\nВ нем представлен Учет личных трат.")
}

func ping(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintln(w, "Так нельзя. Код: 405")
	} else {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "pong. Код: 200")
	}
}

func main() {
	db, err := sql.Open("sqlite", "expenses.db?_pragma=foreign_keys(1)")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	schema, err := os.ReadFile("migrations/create_table.sql")
	if err != nil {
		log.Fatal(err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewExpenseRepository(db)
	handlers.Repo = repo
	handlers.Users = repository.NewUserRepository(db)
	handlers.Sessions = repository.NewSessionRepository(db)
	http.HandleFunc("/", sait)
	http.HandleFunc("/about", about)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("GET /expenses", handlers.RequireAuth(handlers.List))
	http.HandleFunc("GET /expenses/new", handlers.RequireAuth(handlers.NewForm))
	http.HandleFunc("POST /expenses/new", handlers.RequireAuth(handlers.Create))
	http.HandleFunc("GET /expenses/edit", handlers.RequireAuth(handlers.EditForm))
	http.HandleFunc("POST /expenses/edit", handlers.RequireAuth(handlers.Update))
	http.HandleFunc("POST /expenses/delete", handlers.RequireAuth(handlers.Delete))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("GET /register", handlers.RegisterForm)
	http.HandleFunc("POST /register", handlers.Register)
	http.HandleFunc("GET /login", handlers.LoginForm)
	http.HandleFunc("POST /login", handlers.Login)
	http.HandleFunc("POST /logout", handlers.Logout)
	fmt.Println("Сервер запущен!")
	http.ListenAndServe(":8080", nil)
}
