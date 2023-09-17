package app

import (
	"fmt"
	"time"
)

type ParseResult struct {
	id   uint8
	html string
}

func UpdateContent(content *string) {
	const n uint8 = 4
	ch := make(chan ParseResult, n)

	go parseGithub(0, ch, "?since=daily")
	// go parseGithub(1, ch, "?since=weekly")
	go parseHackerNews(1, ch, "/newest")
	go parseHackerNews(2, ch, "/")
	go parseHackerNews(3, ch, "/show")
	// go parseYandexPages(4, ch)
	// go parseHabrPages(4, ch)

	// await for all results, and sort them by id
	parseResults := [n]ParseResult{}
	for _, res := range [n]ParseResult{<-ch, <-ch, <-ch, <-ch} {
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
