package main

import (
	"fmt"
	"myserver/internal/handlers"
	"net/http"
)

func sait(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, это главная страница!")
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Это мой первый проект на языке программирования GO!")
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
	http.HandleFunc("/", sait)
	http.HandleFunc("/about", about)
	http.HandleFunc("/ping", ping)
	http.HandleFunc("/expenses", handlers.List)
	http.HandleFunc("GET /expenses/new", handlers.NewForm)
	http.HandleFunc("POST /expenses/new", handlers.Create)
	fmt.Println("Сервер запущен!")
	http.ListenAndServe(":8080", nil)
}
