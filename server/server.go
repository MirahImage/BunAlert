package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"

	bun "github.com/MirahImage/BunAlert/bun"
)

const port = 9090

const httpContentTypeValue = "text/html; charset=utf-8"
const httpContentTypeHeader = "Content-Type"

const indexTemplate = "index.gtpl"
const reportTemplate = "report.gtpl"
const successTemplate = "success.gtpl"

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome")

	w.Header().Set(httpContentTypeHeader, httpContentTypeValue)

	err := writeTemplate(w, indexTemplate, nil)
	if err != nil {
		log.Fatal("Write template failed:", err)
	}
}

func bunReport(w http.ResponseWriter, r *http.Request) {
	log.Println("Bun report")
	log.Println("Method:", r.Method)

	w.Header().Set(httpContentTypeHeader, httpContentTypeValue)

	if r.Method == "GET" {
		err := writeTemplate(w, reportTemplate, nil)
		if err != nil {
			log.Fatal("Write Template failed:", err)
		}
	} else {
		r.ParseForm()
		size, err := strconv.Atoi(template.HTMLEscapeString(r.Form.Get("size")))
		if err != nil {
			size = 0
			log.Println("Expected integer size ", err)
		}
		var b bun.Bun
		b.LogBun(size, template.HTMLEscapeString(r.Form.Get("description")))
		fmt.Println("Bun Size:", b.Size)
		fmt.Println("Bun Description:", b.Description)
		log.Println("Bun reported")
	}
}

func writeTemplate(w http.ResponseWriter, templateFile string, templateData interface{}) error {
	t, err := template.ParseFiles(templateFile)

	if err != nil {
		return err
	}

	err = t.Execute(w, templateData)
	return err
}

func main() {
	listeningAddress := fmt.Sprintf(":%d", port)
	log.Println(fmt.Sprintf("Listening on %s", listeningAddress))

	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/bunReport", bunReport)

	err := http.ListenAndServe(listeningAddress, nil)

	if err != nil {
		log.Fatal("Listen and Serve:", err)
	}
}
