package main

import (
	"fmt"
	"net/http"
)

func sait(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Привет, это главная страница!")	
}

func about(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Это мой первый проект на языке программирования GO!")
}

func ping(w http.ResponseWriter, r *http.Request){
	if  r.Method != "GET" {
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

	fmt.Println("Сервер запущен!")
	http.ListenAndServe(":8080", nil)
}