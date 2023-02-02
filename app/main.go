package main

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

type parseRes struct {
	id   int
	html string
}

func updateHtml() {
	channel := make(chan parseRes, 3)

	// parse websites at the same time
	go parseGithub(0, channel)
	go parseHackerNews(1, channel, "/")
	go parseHackerNews(2, channel, "/show")

	// await for all parse websites results, and sort them by id
	parseResults := make([]parseRes, 3, 3)
	for _, res := range []parseRes{<-channel, <-channel, <-channel} {
		parseResults[res.id] = res
	}

	// add all results to body
	body := ""
	for _, res := range parseResults {
		body += `<div class="column">` + res.html + `</div>`
	}

	// add body to base html, then write file
	file, _ := os.Create("./public/index.html")
	file.WriteString(fmt.Sprintf(baseHtml, body))
	defer file.Close()

	fmt.Println("-- updated", time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * 60 * 1)
	updateHtml()
}

func main() {
	go updateHtml()

	// start server
	http.Handle("/", http.FileServer(http.Dir("./public")))
	http.ListenAndServe(":3000", nil)
}
