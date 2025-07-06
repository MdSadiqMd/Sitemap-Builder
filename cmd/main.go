package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"os"

	"github.com/MdSadiqMd/Sitemap-Builder/pkg"
)

const xmlns = "https://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	urlFlag := flag.String("url", "https://google.com", "The url for which we want to build sitemap for")
	maxDepth := flag.Int("depth", 2, "Depth that can be crawled upto")
	flag.Parse()

	pages := pkg.BFS(*urlFlag, *maxDepth)
	toXml := urlset{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", " ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}
	fmt.Println()
}
