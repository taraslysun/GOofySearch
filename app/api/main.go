package main

import (
	setup "back/search"
	"back/server"
	"fmt"
)

func main() {

	es := setup.CreateClient()
	fmt.Println("Server is running")
	server.Run(es)

}
