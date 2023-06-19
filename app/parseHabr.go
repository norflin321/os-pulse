package app

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

type HabrItem struct {
	title, link, info string
}

func parseHabr(url string, ch chan []HabrItem) {
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
	ch <- items
}

func parseHabrPages(id int, channel chan ParseResult) {
	fmt.Println("Parse Habr...")
	ch := make(chan []HabrItem, 2)

	go parseHabr("https://habr.com/ru/all/top10/", ch)
	go parseHabr("https://habr.com/ru/all/top10/page2", ch)

	items := []HabrItem{}
	items = append(items, <-ch...)
	items = append(items, <-ch...)
	items = unique(items)

	// create html
	itemsHtml := ""
	for _, item := range items {
		itemsHtml += sprintfSafely(defaultItemHtmlTemplate, item.link, item.title, item.info)
	}

	channel <- ParseResult{id, fmt.Sprintf(columnHtml, itemsHtml)}
}
