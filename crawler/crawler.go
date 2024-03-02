package main

import (
	"os"
	"strconv"

	scraper "example.com/scraper"
)
func main() {

	var visitedLinks = make(map[string]bool)
	var pageFile, _ = os.Create(os.Args[1])
	url := os.Args[2]
	var i,_ = strconv.Atoi(os.Args[3])
	scraper.Scraper(&url, &visitedLinks, i, pageFile)
}
