package scraper

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gocolly/colly/v2"
)


func Scraper(url *string, visitedLinks *map[string]bool, i int, pageFile *os.File) {
	c := colly.NewCollector()

	var links []string
	var pageText string

	visitLink := func(url *string) {
		if !(*visitedLinks)[*url] {
			err := c.Visit(*url)
			if err != nil {
				log.Println("Error visiting link:", err)
				return
			}
		} else {
			// fmt.Println("Link already visited:", *url)
		}
		Scraper(url, visitedLinks, i, pageFile)
	}

	c.OnHTML("a", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		link = strings.Split(link, "#")[0]
		if !strings.Contains(link, "http") {
			domain := strings.Split(*url, "/")
			link = domain[0] + "//" + domain[2] + link
		}
		if !strings.Contains(link, "javascript") && !(*visitedLinks)[link] {
			(*visitedLinks)[link] = true
			links = append(links, link)
			i--
			if i > 0{
				fmt.Println(len((*visitedLinks)), "\tLink found:", link)
				fmt.Println("Visiting link:", link)
				visitLink(&link)
			} else {
				fmt.Println("Done")
				pageFile.Close()
				os.Exit(0)
			}
		}
	})

	c.OnResponse(func(r *colly.Response) {
		pageText = string(r.Body)
		startTitle := strings.Index(pageText, "<title>") + 7
		var endTitle int
		if startTitle == -1 {
			startTitle = 0
			endTitle = 100
		} else {
			endTitle = strings.Index(pageText, "</title>")
			if endTitle == -1 {
				endTitle = 100
			}
		}
		pageTitle := pageText[startTitle:endTitle]
		pageText = strings.ReplaceAll(pageText, "\n", " ")
		pageText = strings.ReplaceAll(pageText, "\"", "'")
		pageTitle = strings.ReplaceAll(pageTitle, "\n", " ")
		pageTitle = strings.ReplaceAll(pageTitle, "\"", "'")
		_, err := pageFile.WriteString(fmt.Sprintf("{\"title\": \"%s\", \"url\": \"%s\", \"text\": \"%s\"},\n", pageTitle, r.Request.URL, pageText))
		if err != nil {
			log.Fatal(err)
		}
	})

	c.OnError(func(res *colly.Response, err error) {
		fmt.Println("OnError: ", err)
	})

	c.Visit(*url)

}
