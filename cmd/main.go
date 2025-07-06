package main

import (
	"flag"
	"fmt"

	"github.com/MdSadiqMd/Sitemap-Builder/pkg"
)

func main() {
	urlFlag := flag.String("url", "https://google.com", "The url for which we want to build sitemap for")
	maxDepth := flag.Int("depth", 3, "Depth that can be crawled upto")
	flag.Parse()

	pages := pkg.BFS(*urlFlag, *maxDepth)
	for _, page := range pages {
		fmt.Println(page)
	}
}
