package pkg

import "strings"

func Filter(base string, links []string) []string {
	var ret []string
	for _, link := range links {
		if strings.HasPrefix(link, base) {
			ret = append(ret, link)
		}
	}
	return ret
}
