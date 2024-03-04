package task_manager

import (
	"fmt"
	"webcrawler/crawler"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gocolly/colly"
)

func Run(c *colly.Collector, url string, es *elasticsearch.Client) {

	links := crawler.Crawl(c, url, es)
	fmt.Println(len(links))

}
