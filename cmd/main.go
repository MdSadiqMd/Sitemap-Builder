package main

import (
	"flag"
	"os"

	handler "github.com/MdSadiqMd/Sitemap-Builder/internal"
)

func main() {
	urlFlag := flag.String("url", "https://google.com", "The url for which we want to build sitemap for")
	maxDepth := flag.Int("depth", 2, "Depth that can be crawled upto")
	flag.Parse()

	err := handler.GenerateSitemap(*urlFlag, *maxDepth, os.Stdout)
	if err != nil {
		panic(err)
	}
}
