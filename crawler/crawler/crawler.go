package crawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"

	"github.com/gocolly/colly"
)

type Result struct {
	Title string `json:"title"`
	Text  string `json:"text"`
	Url   string `json:"url"`
}

func Index(es *elasticsearch.Client, r *Result) {
	data, err := json.Marshal(r)
	if err != nil {
		fmt.Println("error")
	}
	fmt.Println("Indexed url: ", r.Url)
	reader := bytes.NewReader(data)

	_, err = es.Index("test", reader)
	if err != nil {
		return
	}
}

func Crawl(c *colly.Collector, url *string, es *elasticsearch.Client) []string {

	var links []string
	var lang = ""
	r := Result{"", "", ""}
	r.Url = *url

	c.OnHTML("html", func(e *colly.HTMLElement) {

		lang = e.Attr("lang")

		if lang == "en" || lang == "en-US" {
			e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
				link := el.Attr("href")
				link = strings.Split(link, "#")[0]
				if !strings.Contains(link, "http") {
					domain := strings.Split(*url, "/")
					link = domain[0] + "//" + domain[2] + link
				}
				links = append(links, link)
			})
		}

	})

	c.OnHTML("head title", func(e *colly.HTMLElement) {
		title := e.Text
		title = strings.ReplaceAll(title, "\n", " ")
		title = strings.ReplaceAll(title, "\"", "'")
		r.Title = title
	})

	c.OnHTML("html body", func(e *colly.HTMLElement) {

		html := e.DOM
		html.Find("script").Remove()
		html.Find("style").Remove()
		html.Find("img").Remove()

		text := html.Text()
		text = strings.ReplaceAll(text, "\n", " ")
		text = strings.ReplaceAll(text, "\t", " ")
		r.Text = text
	})

	c.OnError(func(res *colly.Response, err error) {
		fmt.Println("OnError: ", err)
	})

	err := c.Visit(*url)
	if err != nil {
		return nil
	}
	if lang == "en" || lang == "en-US" {
		Index(es, &r)
	}
	return links

}
