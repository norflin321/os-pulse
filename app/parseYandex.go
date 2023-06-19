package app

import (
	"fmt"
	"sort"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type YandexItem struct {
	title, link, info string
	timestamp         int64
}

func parseYandexAcademyPage(pageChannel chan []YandexItem) {
	fmt.Println("Parse Yandex Academy...")
	var postClassNames = [...]string{".PostPreviewDouble_preview-double__H_5Se", ".PostPreviewTriple_preview-triple__bYk8m", ".PostPreviewMinor_preview-minor__dc_bH", ".PostPreviewBanner_preview-banner__zN64X", ".PostPreviewYellow_preview-yellow__6oVVr"}
	const url = "https://academy.yandex.ru/journal"
	items := []YandexItem{}

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

				// get date
				infoSlice := strings.Split(item.info, " ")
				if len(infoSlice) >= 3 {
					ruDate := strings.Join(infoSlice[len(infoSlice)-3:], " ")
					date, err := parceRuDate(ruDate)
					if err == nil {
						item.timestamp = date.Unix()
					}
				}
				if item.timestamp == 0 {
					return
				}

				items = append(items, item)
			})
		}
	})
	pageChannel <- items
}

func parseYandexCodePage(pageChannel chan []YandexItem) {
	fmt.Println("Parse Yandex Code...")
	const url = "https://thecode.media/"
	items := []YandexItem{}

	fetchHtmlPage(url).Find(".main-category-post").Each(func(_ int, postEl *goquery.Selection) {
		item := YandexItem{}

		item.title = PrettyStr(postEl.Find(".main-category-post__title").Text())
		item.link, _ = postEl.Attr("href")
		item.info = PrettyStr(postEl.Find(".main-category-post__date").Text())

		// get date
		if item.info != "" {
			date, err := parceRuDate(item.info)
			if err == nil {
				item.timestamp = date.Unix()
			}
		}
		if item.timestamp == 0 {
			return
		}

		items = append(items, item)
	})
	pageChannel <- items
}

func parseYandexPages(id int, channel chan ParseResult) {
	pageChannel := make(chan []YandexItem, 2)
	go parseYandexAcademyPage(pageChannel)
	go parseYandexCodePage(pageChannel)

	items := []YandexItem{}
	items = append(items, <-pageChannel...)
	items = append(items, <-pageChannel...)
	items = unique(items)

	// sort by time of publication
	sort.Slice(items, func(i, j int) bool {
		return items[i].timestamp > items[j].timestamp
	})

	// create html
	itemsHtml := ""
	for _, item := range items {
		itemsHtml += sprintfSafely(defaultItemHtmlTemplate, item.link, item.title, item.info)
	}
	channel <- ParseResult{id, fmt.Sprintf(columnHtml, itemsHtml)}
}
