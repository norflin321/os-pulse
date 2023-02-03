package main

import (
	"embed"
	"html/template"
	"log"
	"net/http"

	"github.com/norflin321/os-pulse/app"
)

//go:embed static/*
var static embed.FS

var html = template.Must(template.ParseFS(static, "static/index.html.tmpl"))
var css, _ = static.ReadFile("static/index.css")

var content string

func main() {
	go app.UpdateContent(&content)

	http.HandleFunc("/", func(res http.ResponseWriter, req *http.Request) {
		data := map[string]interface{}{
			"content": template.HTML(content),
			"css":     template.CSS(app.PrettyStr(string(css))),
		}
		html.ExecuteTemplate(res, "index.html.tmpl", data)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
