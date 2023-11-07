package main

import (
	"html/template"
	"log"
	"net/http"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("articles").Parse(articlesTemplate)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		http.Error(w, err.Error(), http.StatusMethodNotAllowed)
		return
	}

	err = t.Execute(w, wb.GetAllArticles())
	if err != nil {
		log.Println(err)
		return
	}

}
