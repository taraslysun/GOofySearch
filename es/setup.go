package es

import (
	"github.com/elastic/go-elasticsearch/v7"
	// utils "webcrawler/utils"
)

func CreateClient() *elasticsearch.Client {
	/*
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
	*/
	client, err := elasticsearch.NewDefaultClient()
	if err != nil {
		return nil
	}
	return client
}
