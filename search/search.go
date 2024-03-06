package search

import (
	"back/utils"
	"context"
	"encoding/json"
	"log"

	"github.com/elastic/go-elasticsearch/v8"
)

func Search(query string, es *elasticsearch.Client) utils.Response {

	searchResp, err := es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("test"),
		es.Search.WithQuery(query),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	if err != nil {
		log.Fatalf("Error getting response: %s", err)
	}

	var r utils.Response

	if err := json.NewDecoder(searchResp.Body).Decode(&r); err != nil {
		log.Fatalf("Error parsing the response body: %s", err)
	}

	return r

}
