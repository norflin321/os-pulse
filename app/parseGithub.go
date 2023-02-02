package main

import (
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type GithubItem struct {
	link, title, desc, langColor, lang, stars, forks, starsToday string
}

func parseGithub(id int, channel chan parseRes) {
	fmt.Println("Parse Github...")
	const url = "https://github.com/trending"
	items := make([]GithubItem, 0)

	// parse wepage and collect information
	fetchHtmlPage(url).Find("article").Each(func(_ int, article *goquery.Selection) {
		item := GithubItem{}
		link, _ := article.Find(".lh-condensed").Find("a").Attr("href")
		item.link = "https://github.com" + link
		item.title = prettyStr(article.Find("h1").Text())
		item.desc = prettyStr(article.Find("p").Text())

		// find language element and color
		langColorEl := article.Find(".repo-language-color").Nodes
		langColor := ""
		if len(langColorEl) > 0 {
			colorAttr := langColorEl[0].Attr[1].Val
			langColor = colorAttr[len(colorAttr)-7:]
		}
		item.langColor = langColor

		// find info element
		info := prettyStr(article.Find(".d-inline-block").Text())
		infoSlice := strings.Fields(info)
		if len(langColor) == 0 {
			infoSlice = append([]string{""}, infoSlice...)
		}

		item.lang = infoSlice[0]
		item.stars = infoSlice[1]
		item.forks = infoSlice[2]
		item.starsToday = infoSlice[5]

		items = append(items, item)
	})

	// create html
	itemsHtml := ""
	for _, item := range items {
		langDiv := ""
		if item.lang != "" {
			langDiv = `
			<div class="lang">
				<div class="icon" style="background-color: %s"></div>
				<div class="text">%s</div>
			</div>
			`
			langDiv = fmt.Sprintf(langDiv, item.langColor, item.lang)
		}
		itemsHtml += fmt.Sprintf(githubItemHtml, item.link, item.title, item.desc, langDiv, item.stars, item.forks, item.starsToday)
	}

	channel <- parseRes{id, fmt.Sprintf(columnHtml, "Github Trending", url, itemsHtml)}
}
