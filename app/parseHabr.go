package app

import (
	"fmt"
	"sort"

	"github.com/PuerkitoBio/goquery"
)

type HabrItem struct {
	title, link, info string
}

type ParseHabrPageRes struct {
	id    uint8
	items []HabrItem
}

func parseHabr(id uint8, url string, ch chan ParseHabrPageRes) {
	items := []HabrItem{}
	fetchHtmlPage(url).Find(".tm-articles-list__item").Each(func(_ int, el *goquery.Selection) {
		item := HabrItem{}

		item.title = el.Find(".tm-title").Text()
		link, _ := el.Find(".tm-title__link").Attr("href")
		item.link = "https://habr.com" + link
		time := el.Find("time").Text()
		rate := el.Find(".tm-votes-meter__value").Text()
		item.info = fmt.Sprintf("%s рейтинг %s", rate, time)

		items = append(items, item)
	})

	ch <- ParseHabrPageRes{id: id, items: items}
}

func parseHabrPages(id uint8, channel chan ParseResult) {
	fmt.Println("Parse Habr...")
	ch := make(chan ParseHabrPageRes, 2)

	go parseHabr(1, "https://habr.com/ru/all/top10/", ch)
	go parseHabr(2, "https://habr.com/ru/all/top10/page2", ch)

	// await until all goroutines finish, then sort results by id
	pagesResults := []ParseHabrPageRes{<-ch, <-ch}
	sort.Slice(pagesResults, func(i, j int) bool {
		return pagesResults[i].id < pagesResults[j].id
	})

	// create html
	itemsHtml := ""
	for _, result := range pagesResults {
		for _, item := range result.items {
			itemsHtml += sprintfSafely(defaultItemHtmlTemplate, item.link, item.title, item.info)
		}
	}

	channel <- ParseResult{id, fmt.Sprintf(columnHtml, itemsHtml)}
}
