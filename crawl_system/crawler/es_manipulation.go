package crawler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

var CFG = elasticsearch.Config{
	CloudID: "mini-google:ZXVyb3BlLXdlc3Q5LmdjcC5lbGFzdGljLWNsb3VkLmNvbTo0NDMkZWRhMWY0MTkyZmJiNGM3YjhiNDQ2ODk4NjBiNGMyNTckOTUzOThlMjVmNjdmNDA4MzhiYzJhOTE4ODAyZDZjYmI=",
	APIKey:  "NFBZcGFJOEJ1WDg3RXdUSUlaX2o6M0hsZkdsWGlSeEtZc1M0NGpqUXkzZw==",
}

func Setup() *elasticsearch.Client {
	es, err := elasticsearch.NewClient(CFG)
	if err != nil {
		log.Fatalf("Error creating the client: %s", err)
	}

	infores, err := es.Info()
	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	fmt.Println(infores)

	return es
}

func IndexData(title, pageText, link string, es *elasticsearch.Client) {
	doc := map[string]interface{}{
		"title": title,
		"text":  pageText,
		"link":  link,
	}

	body, err := json.Marshal(doc)
	if err != nil {
		fmt.Println("Error marshalling document:", err)
		return
	}

	req := esapi.IndexRequest{
		Index:      "crawl_data",
		DocumentID: strconv.Itoa(id),
		Body:       bytes.NewReader(body),
		Refresh:    "true",
	}
	id++

	res, err := req.Do(context.Background(), es)
	if err != nil {
		fmt.Println("Error indexing document:", err)
		return
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)

	if res.IsError() {
		fmt.Println("Failed to index document:", res.Status())
		return
	}

}
