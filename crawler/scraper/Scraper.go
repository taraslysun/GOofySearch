package scraper

import (
	"fmt"
	"strings"
	"log"
	"os"

	"github.com/gocolly/colly/v2"
)

var visitedLinks = make(map[string]bool)

func Scraper(url *string) {
	c := colly.NewCollector()

	var links []string
	var pageText string

	visitLink := func(url string) {
		err := c.Visit(url)
		if err != nil {
			// log.Println("Error visiting link:", err)
			return
		}
	}

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		if !strings.Contains(link, "javascript") && !visitedLinks[link] {
			visitedLinks[link] = true
			links = append(links, link)
			fmt.Println(len(visitedLinks), "Link found:", link)
			visitLink(link)
		}
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {
		pageText += e.Text + "\n"
	})

	c.OnError(func(res *colly.Response, err error) {
		// log.Println("Something went wrong:", err)
	})

	err := c.Visit(*url)
	if err != nil {
		// log.Fatal(err)
	}

	pageFile, err := os.Create("page_text.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer pageFile.Close()
	_, err = pageFile.WriteString(pageText)
	if err != nil {
		log.Fatal(err)
	}

	linksFile, err := os.Create("links.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer linksFile.Close()
	for _, link := range links {
		_, err := linksFile.WriteString(link + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

}
