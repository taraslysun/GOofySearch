package main

import (
	// "back/search"
	setup "back/search"
	"back/server"
)

func main() {

	es := setup.CreateClient()
	// search.Ingest(es)
	server.Run(es)

}
