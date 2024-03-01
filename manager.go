package manager

import (
	"fmt"
	"github.com/gocolly/colly"
)

func manager() {
	c := colly.NewCollector()
	fmt.Println(c)
}
