package main

import (
	manager "webcrawler/task_manager"

	setup "webcrawler/es"

	"github.com/gocolly/colly"
)

func main() {

	es := setup.CreateClient()
	c := colly.NewCollector()
	manager.Run(c, "https://en.wikipedia.org/wiki/Main_Page", es)

}
