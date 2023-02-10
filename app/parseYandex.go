package app

import (
	"fmt"

	"github.com/PuerkitoBio/goquery"
)

var postClassNames = [...]string{
	".PostPreviewDouble_preview-double__H_5Se",
	".PostPreviewTriple_preview-triple__bYk8m",
	".PostPreviewMinor_preview-minor__dc_bH",
	".PostPreviewBanner_preview-banner__zN64X",
	".PostPreviewYellow_preview-yellow__6oVVr",
}

type YandexItem struct {
	title, link, info string
}

func parseYandex(id int, channel chan ParseResult) {
	fmt.Println("Parse Yandex...")
	const url = "https://academy.yandex.ru/journal"
	items := make([]YandexItem, 0)

	// try to find rows of posts
	fetchHtmlPage(url).Find(".Pattern_pattern__57Vkx").Each(func(_ int, rowEl *goquery.Selection) {
		// in each row, try to find posts (each post can have different className)
		for _, className := range postClassNames {
			rowEl.Find(className).Each(func(_ int, postEl *goquery.Selection) {
				item := YandexItem{}
				// in each post, try to find links
				postEl.Find("a").Each(func(i int, linkEl *goquery.Selection) {
					text := linkEl.Text()
					href, _ := linkEl.Attr("href")
					if i == 0 {
						item.title = text
						item.link = "https://academy.yandex.ru" + href
					}
					if i == 1 && text != "узнать больше" {
						item.info = text + " "
					}
				})
				timeEl := postEl.Find("time")
				item.info += timeEl.Text()
				items = append(items, item)
			})
		}
	})

	// create html
	itemsHtml := ""
	for _, item := range items {
		itemsHtml += fmt.Sprintf(yandexItemHtml, item.link, item.title, item.info)
	}
	channel <- ParseResult{id, fmt.Sprintf(columnHtml, "Yandex Academy", url, itemsHtml)}
}
