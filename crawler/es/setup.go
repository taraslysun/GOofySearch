package es

import (
	"fmt"
	"log"
	utils "webcrawler/utils"

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

