package main

import (
	"github.com/gocolly/colly"
	setup "webcrawler/es"
	manager "webcrawler/task_manager"
)

func main() {
	es := setup.CreateClient()
	c := colly.NewCollector()
	manager.Run(c, "https://www.wikipedia.org/", es)
}
