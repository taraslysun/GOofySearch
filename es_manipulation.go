package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v7"
)

func Setup() *elasticsearch.Client {
	// create a new Elasticsearch client

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		fmt.Println("Error creating Elasticsearch client:", err)
		return nil
	}

	return es
}

func IndexData(title, pageText, link string, es *elasticsearch.Client) {
	// index data to elasticsearch
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
	defer res.Body.Close()

	if res.IsError() {
		fmt.Println("Failed to index document:", res.Status())
		return
	}

}
