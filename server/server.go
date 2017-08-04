package main

import (
	"fmt"
	"net/http"
	"log"
	"html/template"
)

const port = 9090

const httpContentTypeValue = "text/html; charset=utf-8"
const httpContentTypeHeader = "Content-Type"


func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome")

	w.Header().Set(httpContentTypeHeader, httpContentTypeValue)

	t, err := template.ParseFiles("index.gtpl")
	if err != nil {
		log.Fatal("template.ParseFiles:", err)
	}
	err = t.Execute(w, nil)
	if err != nil {
		log.Fatal("template.Execute:", err)
	}
}

func bunReport(w http.ResponseWriter, r *http.Request) {
	log.Println("Bun report")
	log.Println("Method:", r.Method)
	
	w.Header().Set(httpContentTypeHeader, httpContentTypeValue)

	if r.Method == "GET" {
		t, err := template.ParseFiles("report.gtpl")
		if err != nil {
			log.Fatal("template.ParseFiles:", err)
		}
		err = t.Execute(w, nil)
		if err != nil {
			log.Fatal("template.Execute:", err)
		}
	} else {
		r.ParseForm()
		fmt.Println("Bun Size:", r.Form.Get("size"))
		fmt.Println("Bun Description:", template.HTMLEscapeString(r.Form.Get("description")))
	}
}

func main() {
	listeningAddress := fmt.Sprintf(":%d", port)
	log.Println(fmt.Sprintf("Listening on %s", listeningAddress))
	
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/bunReport", bunReport)
	
	err := http.ListenAndServe(":9090", nil)

	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
