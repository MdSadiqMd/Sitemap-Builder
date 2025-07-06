package pkg

func BFS(urlStr string, maxDepth int) []string {
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: {},
	}

	for range maxDepth {
		q, nq = nq, make(map[string]struct{})
		for url := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range Get(url) {
				nq[link] = struct{}{}
			}
		}
	}

	ret := make([]string, 0)
	for url := range seen {
		ret = append(ret, url)
	}
	return ret
}
