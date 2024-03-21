package main

import (
	setup "webcrawler/es"
	manager "webcrawler/task_manager"

	"github.com/gocolly/colly"
)

func main() {
	es := setup.CreateClient()
	_, err := es.Indices.Delete([]string{"test"})
	if err != nil {
		return
	}
	_, err = es.Indices.Create("test")
	if err != nil {
		return
	}
	c := colly.NewCollector()
	manager.Run(c, "https://www.wikipedia.org/", es)
}
