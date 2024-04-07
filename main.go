package main

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/elastic/go-elasticsearch/v7"
	"golang.org/x/net/html"
)

var id = 1

func LinkToChannel(link string, crawledLinksChannel chan string) {
	crawledLinksChannel <- link
}

// MonitorCrawling ends crawling if there is no links to scrape especially needed when working without task manager
func MonitorCrawling(pendingLinksChannel chan string, crawledLinksChannel chan string, linksAmountChannel chan int) {
	i := 0
	for j := range linksAmountChannel {
		i += j

		// check if number of pending links is 0
		// if yes, close all the channels
		if i == 0 {
			close(pendingLinksChannel)
			close(crawledLinksChannel)
			close(linksAmountChannel)
		}
	}
}

// ProcessCrawledLinks used for filtering visited links
func ProcessCrawledLinks(pendingLinksChannel chan string, crawledLinksChannel chan string, linksAmountChannel chan int) {
	foundUrls := make(map[string]bool)

	// iterating over crawled links
	// check if visited ? skip : add to pending links
	for cl := range crawledLinksChannel {
		if !foundUrls[cl] {
			foundUrls[cl] = true
			linksAmountChannel <- 1
			pendingLinksChannel <- cl
		}
	}
}

func getResponse(link *string, agent string) *http.Response {
	req, err := http.NewRequest("GET", *link, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", agent)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error while connecting to site.")
	}
	return resp
}

func formatHref(base, href string) (string, bool) {
	if strings.HasPrefix(href, "http") {
		return href, true
	} else if strings.HasPrefix(href, "/") {
		return base + href, true
	}
	return "", false
}

func extractLink(token *html.Token, link *string) (string, bool) {
	href := ""
	for _, a := range token.Attr {
		if a.Key == "href" {
			href = a.Val
			formattedHref, ok := formatHref(*link, href)
			if ok {
				return formattedHref, true
			}

		}
	}
	return "", false
}

func RandomString(userAgentList []string) string {
	randomIndex := rand.Intn(len(userAgentList))
	return userAgentList[randomIndex]
}

func extractContent(link *string, crawledLinksChannel chan string, es *elasticsearch.Client) {
	userAgentList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18363",
	}
	var title string
	var pageText string

	response := getResponse(link, RandomString(userAgentList))

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(response.Body)

	z := html.NewTokenizer(response.Body)

	for {
		tokenType := z.Next()
		switch tokenType {
		case html.StartTagToken, html.SelfClosingTagToken:
			token := z.Token()

			switch token.Data {

			case "html":
				for _, attr := range token.Attr {
					if attr.Key == "lang" {
						lang := attr.Val
						if lang != "en" && lang != "en-us" && lang != "ua" {
							break
						}

					}
				}

			case "title":
				tokenType = z.Next()
				if tokenType == html.TextToken {
					title = z.Token().Data
				}

			case "h1", "h2", "h3", "h4", "h5", "h6", "span", "p":
				tokenType = z.Next()
				if tokenType == html.TextToken {
					pageText += z.Token().Data
				}

			case "a":
				href, ok := extractLink(&token, link)
				if ok {
					go LinkToChannel(href, crawledLinksChannel)
				}
			}
		default:
			panic("unhandled default case")

		}

		if tokenType == html.ErrorToken {
			break
		}
	}
	if title != "" && pageText != "" {
		fmt.Println(*link)
		IndexData(title, pageText, *link, es)
	}
}

func CrawlWebpage(wg *sync.WaitGroup, pendingLinksChannel chan string,
	crawledLinksChannel chan string, linksAmountChannel chan int, depth int, es *elasticsearch.Client) {

	for link := range pendingLinksChannel {
		extractContent(&link, crawledLinksChannel, es)
		linksAmountChannel <- -1
		depth--
		if depth == 0 {
			break
		}
	}
	wg.Done()
}

func CrawlerMain(startLinks []string, depth int, numThreads int, es *elasticsearch.Client) {
	pendingLinksChannel := make(chan string)
	crawledLinksChannel := make(chan string)
	linksAmountChannel := make(chan int)

	for _, startLink := range startLinks {
		go LinkToChannel(startLink, crawledLinksChannel)
	}

	go ProcessCrawledLinks(pendingLinksChannel, crawledLinksChannel, linksAmountChannel)
	go MonitorCrawling(pendingLinksChannel, crawledLinksChannel, linksAmountChannel)

	var wg sync.WaitGroup

	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go CrawlWebpage(&wg, pendingLinksChannel, crawledLinksChannel, linksAmountChannel, depth, es)
	}

	wg.Wait()

	// post method to task manager to return new links
	// need to return links from pendingLinksChannel

}

func ManageCrawler() {
	/*
		for {
		get request to task manager to get links

		give to crawler 15-20 links with desired depth
		CrawlerMain(link_from_task_manager, desired_depth, number_of_threads)
		if depth is > 1, number_of_threads has to be greater than number of links (desirable)

		check some condition to end cycle
		}
	*/
}

// Run crawler without task manager
func main() {
	es := Setup()
	fmt.Println("Crawl started!...")
	links := []string{"https://en.wikipedia.org/wiki/Nelson-class_battleship"}
	CrawlerMain(links, 2, 8, es)
}
