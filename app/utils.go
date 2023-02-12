package app

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func fetchHtmlPage(url string) *goquery.Document {
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}
	html, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	return html
}

func PrettyStr(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Join(strings.Fields(str), " ")
	str = strings.TrimSpace(str)
	return str
}

func prettyStruct(v interface{}) string {
	return fmt.Sprintf("%+v", v)
}

func prettyMap(a interface{}) string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}

func unique[T comparable](s []T) []T {
	inResult := map[T]bool{}
	var result []T
	for _, str := range s {
		if _, ok := inResult[str]; !ok {
			inResult[str] = true
			result = append(result, str)
		}
	}
	return result
}

func reverseKeyVal(m map[string]string) map[string]string {
	n := make(map[string]string, len(m))
	for k, v := range m {
		n[v] = k
	}
	return n
}

// date example "07 марта 1999"
// golang cheatsheet for dates:
// https://gosamples.dev/date-time-format-cheatsheet/
func parceRuDate(inDate string) (time.Time, error) {
	var ruMonthsMapping = map[string]string{"января": "01", "февраля": "02", "марта": "03", "апреля": "04", "мая": "05", "июня": "06", "июля": "07", "августа": "08", "сентября": "09", "октября": "10", "ноября": "11", "декабря": "12"}

	dateSplit := strings.Split(inDate, " ")
	dayInt, _ := strconv.ParseInt(dateSplit[0], 10, 0)
	dateSplit[0] = fmt.Sprintf("%02d", dayInt)
	dateSplit[1] = ruMonthsMapping[dateSplit[1]]

	outDate, err := time.Parse("02-01-2006", strings.Join(dateSplit, "-"))
	if err != nil {
		return time.Time{}, err
	}
	return outDate, nil
}

func sprintfSafely(outHtml string, insertingStrings ...string) string {
	var escapedStrings []any
	for _, str := range insertingStrings {
		escapedStrings = append(escapedStrings, html.EscapeString(str))
	}
	return fmt.Sprintf(outHtml, escapedStrings...)
}
