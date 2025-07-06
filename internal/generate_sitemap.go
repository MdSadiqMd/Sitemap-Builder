package handler

import (
	"encoding/xml"
	"fmt"
	"io"

	"github.com/MdSadiqMd/Sitemap-Builder/pkg"
)

const xmlns = "https://www.sitemaps.org/schemas/sitemap/0.9"

type Loc struct {
	Value string `xml:"loc"`
}

type UrlSet struct {
	Urls  []Loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func GenerateSitemap(url string, maxDepth int, writer io.Writer) error {
	pages := pkg.BFS(url, maxDepth)
	toXml := UrlSet{
		Xmlns: xmlns,
	}
	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, Loc{page})
	}

	fmt.Fprint(writer, xml.Header)
	enc := xml.NewEncoder(writer)
	enc.Indent("", " ")
	if err := enc.Encode(toXml); err != nil {
		return err
	}
	fmt.Fprintln(writer)
	return nil
}
