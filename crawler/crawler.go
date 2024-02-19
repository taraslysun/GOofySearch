package main

import (
	"example.com/scraper"
)

func main() {

	url := "https://cms.ucu.edu.ua/?redirect=0"
	scraper.Scraper(&url)
}
