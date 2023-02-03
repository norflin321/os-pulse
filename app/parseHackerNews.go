package app

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HNItem struct {
	sitebit, link, title, info string
}

func parseHackerNews(id int, channel chan ParseResult, page string) {
	fmt.Println("Parse HackerNews...", page)
	const url = "https://news.ycombinator.com"
	items := make([]HNItem, 0)

	html := fetchHtmlPage(url + page)

	// collect information from title element
	html.Find(".athing").Each(func(i int, s *goquery.Selection) {
		item := HNItem{}

		titleEl := s.Find(".titleline")
		link, _ := titleEl.Find("a").Attr("href")
		sitebitEl := titleEl.Find(".sitebit").Find("a")
		sitebit := ""
		if len(sitebitEl.Nodes) > 0 {
			sitebit, _ = sitebitEl.Attr("href")
			sitebit = sitebit[10:]
		}
		item.sitebit = sitebit
		if sitebit == "" {
			item.link = "https://news.ycombinator.com/" + link
		} else {
			item.link = link
		}
		item.title = titleEl.Find("a").Nodes[0].FirstChild.Data

		items = append(items, item)
	})

	// collect information from undertitle element
	html.Find(".subtext").Each(func(i int, s *goquery.Selection) {
		info := PrettyStr(s.Text())
		info = strings.ReplaceAll(info, " | hide | ", " ")
		info = strings.ReplaceAll(info, " | hide", " ")
		info = strings.ReplaceAll(info, " | ", " ")
		items[i].info = info
	})

	// create html
	itemsHtml := ""
	for _, item := range items {
		itemsHtml += fmt.Sprintf(hackerNewsItemHtml, item.link, item.title, item.info)
	}

	title := "Hacker News"
	if page == "/show" {
		title += " Show"
	}

	channel <- ParseResult{id, fmt.Sprintf(columnHtml, title, url+page, itemsHtml)}
}
