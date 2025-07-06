package pkg

import (
	"bytes"
	"io"
	"net/http"
	"net/url"
)

func Get(urlString string) []string {
	resp, err := http.Get(urlString)
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
	pages := Hrefs(bodyReader, base)
	return Filter(pages, withPrefix(base))
}
