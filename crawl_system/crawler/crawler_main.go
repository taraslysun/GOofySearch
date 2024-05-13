package crawler

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	"encoding/json"
	"bytes"

	"github.com/elastic/go-elasticsearch/v7"
	"golang.org/x/net/html"
	"log"
)

type AtomicId struct {
	value int64
}

func (id *AtomicId) Get() int64 {
	return atomic.LoadInt64(&id.value)
}

func (id *AtomicId) Set(val int64) {
	atomic.StoreInt64(&id.value, val)
}

func (id *AtomicId) Increment() int64 {
	return atomic.AddInt64(&id.value, 1)
}

var id AtomicId

func LinkToChannel(link *string, wgLinks *sync.WaitGroup, crawledLinksChannel chan string, linksAmountChannel chan int) {
	crawledLinksChannel <- *link
	linksAmountChannel <- 1
	wgLinks.Done()
}

// MonitorCrawling ends crawling if there is no links to scrape especially needed when working without task manager
func MonitorCrawling(pendingLinksChannel, crawledLinksChannel chan string, linksAmountChannel chan int) {
	i := 0
	for j := range linksAmountChannel {
		i += j

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

	for cl := range crawledLinksChannel {
		if !foundUrls[cl] {
			foundUrls[cl] = true
			pendingLinksChannel <- cl
		}
		linksAmountChannel <- -1
	}
}

func getResponse(link *string, agent string) *http.Response {
	req, err := http.NewRequest("GET", *link, nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("User-Agent", agent)

	client := &http.Client{
		Timeout: 2 * time.Second,
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil
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

func extractContent(link *string, wgLinks *sync.WaitGroup, crawledLinksChannel chan string,
	es *elasticsearch.Client, linksAmountChannel chan int, wgIndex *sync.WaitGroup) {
	userAgentList := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/93.0.4577.82 Safari/537.36",
		"Mozilla/5.0 (iPhone; CPU iPhone OS 14_4_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/14.0.3 Mobile/15E148 Safari/604.1",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/87.0.4280.141 Safari/537.36 Edg/87.0.664.75",
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/70.0.3538.102 Safari/537.36 Edge/18.18363",
	}
	var title string
	var pageText string

	response := getResponse(link, RandomString(userAgentList))

	if response == nil {
		return
	}

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
							return
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
					wgLinks.Add(1)
					go LinkToChannel(&href, wgLinks, crawledLinksChannel, linksAmountChannel)
				}
			}
		}

		if tokenType == html.ErrorToken {
			break
		}
	}
	if title != "" && pageText != "" {
		fmt.Println(strconv.Itoa(int(id.Get())), " Content link", *link)
		wgIndex.Add(1)
		go func() {
			IndexData(title, pageText, *link, es)
		}()
		wgIndex.Done()

	}
}

func CrawlWebpage(wg *sync.WaitGroup, pendingLinksChannel chan string,
	crawledLinksChannel chan string, linksAmountChannel chan int, es *elasticsearch.Client, wgIndex *sync.WaitGroup) {

	var wgLinks sync.WaitGroup

	link := <-pendingLinksChannel
	extractContent(&link, &wgLinks, crawledLinksChannel, es, linksAmountChannel, wgIndex)
	wgLinks.Wait()
	linksAmountChannel <- -1

	wg.Done()
}

func CrawlerMain(startLinks []string, numLinks int, es *elasticsearch.Client, masterIp string) {
	pendingLinksChannel := make(chan string, 100)
	crawledLinksChannel := make(chan string, 100000)
	linksAmountChannel := make(chan int, 100)

	var wgStart sync.WaitGroup
	var wgIndex sync.WaitGroup
	defer wgIndex.Wait()

	for _, startLink := range startLinks {
		wgStart.Add(1)

		go func(link string) {
			defer wgStart.Done()

			pendingLinksChannel <- link
			linksAmountChannel <- 1
		}(startLink)
	}

	wgStart.Wait()

	go MonitorCrawling(pendingLinksChannel, crawledLinksChannel, linksAmountChannel)

	var wg sync.WaitGroup

	for i := 0; i < numLinks; i++ {
		wg.Add(1)
		go CrawlWebpage(&wg, pendingLinksChannel, crawledLinksChannel,
			 linksAmountChannel, es, &wgIndex)
	}
	wg.Wait()

	go ProcessCrawledLinks(pendingLinksChannel, crawledLinksChannel, linksAmountChannel)

	var links []string
	var linksMap = make(map[string]struct{})

	for link := range pendingLinksChannel {
		if _, ok := linksMap[link]; !ok {
			links = append(links, link)
			linksMap[link] = struct{}{} 
		}
	}

	linksStr := strings.Join(links, "~")

	payload := map[string]string{
		"links": linksStr,
	}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	resp, err := http.Post("http://" + masterIp+ ":9092/links", "application/json", bytes.NewBuffer(jsonPayload))
	if err != nil {
		fmt.Println("Error making POST request:", err)
		return
	}
	defer resp.Body.Close()

}


func MasterCrawler(es *elasticsearch.Client, masterIp string) {
	client := &http.Client{}
	var wg sync.WaitGroup

	for {
		for i := 1; i <= 5; i++ {
			wg.Add(1)

			go func(id int) {
				defer wg.Done()

				res, err := http.NewRequest("GET", "http://localhost:8080/links", nil)
				if err != nil {
					log.Fatal(err)
				}

				q := res.URL.Query()
				q.Add("CID", strconv.Itoa(id))
				res.URL.RawQuery = q.Encode()

				resp, err := client.Do(res)
				if err != nil {
					log.Fatal(err)
				}
				defer func(Body io.ReadCloser) {
					err := Body.Close()
					if err != nil {
						log.Fatal(err)
					}
				}(resp.Body)

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatal(err)
				}

				var links []string
				err = json.Unmarshal(body, &links)
				if err != nil {
					log.Fatal(err)
				}

				if len(links) == 0 {
					return
				}

				CrawlerMain(links, len(links), es, masterIp)
			}(1)

		}
		wg.Wait()

	}
}
