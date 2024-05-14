package search

import (
	"back/utils"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func Search(query string, es *elasticsearch.Client) utils.Response {

	var buf bytes.Buffer
	queryBody := map[string]interface{}{
		"query": map[string]interface{}{
			"bool": map[string]interface{}{
				"should": []map[string]interface{}{
					{
						"match": map[string]interface{}{
							"text": map[string]interface{}{
								"query": query,
								"boost": 1,
							},
						},
					},
					{
						"match": map[string]interface{}{
							"title": map[string]interface{}{
								"query": query,
								"boost": 3,
							},
						},
					},
					{
						"match": map[string]interface{}{
							"link": map[string]interface{}{
								"query": query,
								"boost": 3,
							},
						},
					},
				},
			},
		},
		"min_score": 9,
	}
	

	if err := json.NewEncoder(&buf).Encode(queryBody); err != nil {
		log.Fatalf("Error encoding query: %s", err)
	}

	searchResp, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("crawl_data"),
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)
	fmt.Println("Elastic search request")

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	var r utils.Response

	if err := json.NewDecoder(searchResp.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}
	return r

}
