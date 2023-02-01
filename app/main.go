package main

import (
	"fmt"
	"net/http"
	"os"
)

func createIndexHtml() {
	body := ""
	body += parseGithub()

	// write file
	file, _ := os.Create("./public/index.html")
	file.WriteString(fmt.Sprintf(baseHtml, body))
	defer file.Close()
}

func main() {
	go createIndexHtml()

	// start server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":3000", nil)
}
