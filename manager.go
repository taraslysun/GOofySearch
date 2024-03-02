package main

import (
	"fmt"
	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.MaxDepth(1),
	)

	c.OnHTML("div", func(e *colly.HTMLElement) {
		link := e.Attr("class")
		fmt.Println(link)
		err := e.Request.Visit(link)
		if err != nil {
			return
		}
	})

	err := c.Visit("https://www.geeksforgeeks.org/priority-queue-set-1-introduction/")
	if err != nil {
		return
	}
}
