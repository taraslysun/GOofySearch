package main

import (
	setup "back/search"
	"back/server"
)

func main() {

	es := setup.CreateClient()
	server.Run(es)

}
