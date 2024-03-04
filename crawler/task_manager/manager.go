package task_manager

import (
	"fmt"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/gocolly/colly"
	"net/url"
	"time"
	"webcrawler/crawler"
)

type DomainParams struct {
	Links   []string      `json:"links"`
	Visited []string      `json:"visited"`
	Timeout time.Duration `json:"timeout"`
}

func StateInit(links []string) map[string]*DomainParams {
	domains := make(map[string]*DomainParams)
	for _, link := range links {
		parsedURL, err := url.Parse(link)
		if err != nil {
			fmt.Println(err)
			continue
		}
		domain := parsedURL.Hostname()
		domains[domain] = &DomainParams{
			Links:   []string{link},
			Visited: []string{},
			Timeout: time.Second * 10,
		}
	}
	return domains
}

func find(element string, checkLst []string) bool {
	for _, elem := range checkLst {
		if elem == element {
			return true
		}
	}
	return false
}

func Run(c *colly.Collector, urlink string, es *elasticsearch.Client) {
	links := crawler.Crawl(c, urlink, es)
	domains := StateInit(links)
	for i := 0; i < 3; i++ {
		for _, value := range domains {
			var newLink string
			if len(value.Links) > 0 {
				newLink = value.Links[0]
				value.Links = value.Links[1:]
			} else {
				continue
			}
			parsedURL, err := url.Parse(newLink)
			if err != nil {
				fmt.Println(err)
				continue
			}
			newDomain := parsedURL.Hostname()
			if _, ok := domains[newDomain]; !ok {
				domains[newDomain] = &DomainParams{}
			}
			if find(newLink, domains[newDomain].Visited) || find(newLink, domains[newDomain].Links) {
				continue
			} else {
				newLinks := crawler.Crawl(c, newLink, es)
				domains[newDomain].Visited = append(domains[newDomain].Visited, newLink)
				for _, localLink := range newLinks {
					parsedURL, err := url.Parse(localLink)
					if err != nil {
						fmt.Println("Error", err)
						continue
					}
					domain := parsedURL.Hostname()
					if _, ok := domains[domain]; !ok {
						domains[domain] = &DomainParams{}
					}
					if find(localLink, domains[domain].Visited) || find(localLink, domains[domain].Links) {
						continue
					} else {
						if _, ok := domains[domain]; !ok {
							domains[domain] = &DomainParams{}
						}
						domains[domain].Links = append(domains[domain].Links, localLink)
					}
				}
			}
		}
	}
	for i, links := range domains {
		fmt.Println(i, links.Visited)
	}
}
