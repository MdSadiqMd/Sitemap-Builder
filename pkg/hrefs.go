package pkg

import (
	"io"
	"strings"

	href "github.com/MdSadiqMd/Href-Parser/pkg"
)

func Hrefs(r io.Reader, base string) []string {
	links, _ := href.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}
