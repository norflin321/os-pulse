package app

import (
	"fmt"
	"time"
)

type ParseResult struct {
	id   int
	html string
}

func UpdateContent(content *string) {
	channel := make(chan ParseResult, 3)

	// parse websites at the same time
	go parseGithub(0, channel)
	go parseHackerNews(1, channel, "/")
	go parseHackerNews(2, channel, "/show")

	// await for all parse websites results, and sort them by id
	parseResults := make([]ParseResult, 3, 3)
	for _, res := range []ParseResult{<-channel, <-channel, <-channel} {
		parseResults[res.id] = res
	}

	// then add all results to content string
	*content = ""
	for _, res := range parseResults {
		*content += `<div class="column">` + res.html + `</div>`
	}

	fmt.Println("-- updated", time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * 60 * 10)
	UpdateContent(content)
}
