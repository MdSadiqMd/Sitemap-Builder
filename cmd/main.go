package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/MdSadiqMd/Sitemap-Builder/pkg"
)

func main() {
	urlFlag := flag.String("url", "https://google.com", "The url for which we want to build sitemap for")
	flag.Parse()

	resp, err := http.Get(*urlFlag)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	bodyReader := bytes.NewReader(bodyBytes)
	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}

	base := baseUrl.String()
	pages := pkg.Hrefs(bodyReader, base)
	for _, page := range pages {
		fmt.Println(page)
	}
}
