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
	const n = 5
	ch := make(chan ParseResult, n)

	go parseGithub(0, ch)
	go parseHackerNews(1, ch, "/show")
	go parseHackerNews(2, ch, "/")
	go parseYandexPages(3, ch)
	go parseHabrPages(4, ch)

	// await for all results, and sort them by id
	parseResults := [n]ParseResult{}
	for _, res := range [n]ParseResult{<-ch, <-ch, <-ch, <-ch, <-ch} {
		parseResults[res.id] = res
	}

	// then add all results to content string
	*content = ""
	for _, res := range parseResults {
		*content += `<div class="column">` + res.html + `</div>`
	}

	fmt.Println("-- updated", time.Now().Format("2006-01-02 15:04:05"))
	time.Sleep(time.Second * 60 * 5)
	UpdateContent(content)
}
