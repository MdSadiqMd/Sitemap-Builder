package main

import (
	"flag"
	"fmt"

	"github.com/MdSadiqMd/Sitemap-Builder/pkg"
)

func main() {
	urlFlag := flag.String("url", "https://google.com", "The url for which we want to build sitemap for")
	flag.Parse()

	pages := pkg.Get(*urlFlag)
	for _, page := range pages {
		fmt.Println(page)
	}
}
