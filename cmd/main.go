package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"

	href "github.com/MdSadiqMd/Href-Parser/pkg"
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
	links, _ := href.Parse(bodyReader)
	for _, l := range links {
		fmt.Println(l)
	}

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()
	fmt.Println(base)
}
