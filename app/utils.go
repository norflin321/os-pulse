package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

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

func prettyStr(str string) string {
	str = strings.Replace(str, "\t", "", -1)
	str = strings.Replace(str, "\n", "", -1)
	str = strings.Join(strings.Fields(str), " ")
	return str
}

func prettyStruct(v interface{}) string {
	return fmt.Sprintf("%+v \n", v)
}

func prettyMap(a interface{}) string {
	b, _ := json.MarshalIndent(a, "", "  ")
	return string(b)
}
