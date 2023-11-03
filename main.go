package main

import (
	"WildBerries/internal/WildBerries"
	"fmt"
	"log"
	"net/http"
)

func main() {
	wb := WildBerries.NewParser()
	arts := []string{
		"25968980",
		"75563253",
	}
	wb.AddArticles(arts)
	mux := http.NewServeMux()
	mux.HandleFunc("/", HandleIndex)
	fmt.Println("Сервер запущен по адресу http://127.0.0.1:4445")
	log.Fatal("Ошибка: ", http.ListenAndServe(":4445", mux))
}
