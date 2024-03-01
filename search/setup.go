package search

import (
	utils "back/utils"
	"bytes"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func CreateClient() *elasticsearch.Client {
	es, err := elasticsearch.NewClient(utils.CFG)
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

func Ingest(es *elasticsearch.Client) {

	ingestResult, err := es.Bulk(
		bytes.NewReader(utils.BUF.Bytes()),
		es.Bulk.WithIndex("index_name"),
	)

	fmt.Println(ingestResult, err)
}
