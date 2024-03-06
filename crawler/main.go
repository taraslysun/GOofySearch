package main

import (
	setup "webcrawler/es"
	manager "webcrawler/task_manager"

	"github.com/gocolly/colly"
)

func main() {
	es := setup.CreateClient()
	es.Indices.Delete([]string{"test"})
	es.Indices.Create("test")
	c := colly.NewCollector()
	manager.Run(c, "https://www.wikipedia.org/", es)
}
