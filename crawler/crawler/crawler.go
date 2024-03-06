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
	Text string `json:"text"`
	Url string `json:"url"`
}


func Index(es *elasticsearch.Client, r *Result) {
	data, err := json.Marshal(r)
        if err != nil {
            fmt.Println("error")
        }
	
    reader := bytes.NewReader(data)

	es.Index("test", reader)
}


func Crawl(c *colly.Collector, url *string, es *elasticsearch.Client) []string {

	
	var links []string
	r := Result{"", "", ""}
	r.Url = *url


	c.OnHTML("html", func(e *colly.HTMLElement) {
		e.ForEach("a[href]", func(_ int, el *colly.HTMLElement) {
		  link := el.Attr("href")
		  link = strings.Split(link, "#")[0]
		  if !strings.Contains(link, "http") {
			domain := strings.Split(*url, "/")
			link = domain[0] + "//" + domain[2] + link
		  }
		  links = append(links, link)
		})
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

	c.Visit(*url)
	Index(es, &r)

	return links

}
