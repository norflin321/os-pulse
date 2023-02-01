package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type HNItem map[string]string

func parseHackerNews(page string) []HNItem {
	fmt.Println("Parse HackerNews...", page)
	items := make([]HNItem, 0)

	html := fetchHtmlPage("https://news.ycombinator.com" + page)

	// collect information from title element
	html.Find(".athing").Each(func(i int, s *goquery.Selection) {
		item := make(HNItem)

		titleEl := s.Find(".titleline")
		link, _ := titleEl.Find("a").Attr("href")
		sitebitEl := titleEl.Find(".sitebit").Find("a")
		sitebit := ""
		if len(sitebitEl.Nodes) > 0 {
			sitebit, _ = sitebitEl.Attr("href")
			sitebit = sitebit[10:]
		}
		item["sitebit"] = sitebit
		if sitebit == "" {
			item["link"] = "https://news.ycombinator.com/" + link
		} else {
			item["link"] = link
		}
		item["title"] = titleEl.Find("a").Nodes[0].FirstChild.Data

		items = append(items, item)
	})

	// collect information from undertitle element
	html.Find(".subtext").Each(func(i int, s *goquery.Selection) {
		items[i]["info"] = strings.ReplaceAll(prettyStr(s.Text()), " | hide | ", " ")
	})

	return items
}
