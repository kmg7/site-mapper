package sitemapper

import (
	"encoding/xml"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	lp "github.com/kmg7/link-parser"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}
type urlSet struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func MapUrl(target string, depth int) {
	res, err := http.Get(target)
	if err != nil {
		panic(err)

	}
	defer res.Body.Close()
	hostUrl := &url.URL{
		Scheme: res.Request.URL.Scheme,
		Host:   strings.TrimSuffix(res.Request.URL.Host, "/"),
	}
	host := hostUrl.String()
	println(host)

	if res.StatusCode != http.StatusOK {
		println(res.StatusCode)
		return

	}
	siteMap := bfs(host, host, depth)
	xmlData := urlSet{Xmlns: xmlns}
	for _, v := range siteMap {

		xmlData.Urls = append(xmlData.Urls, loc{v})
	}
	fmt.Print(xml.Header)
	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "\t")
	if err := enc.Encode(xmlData); err != nil {
		panic(err)
	}
	fmt.Println()

}

func bfs(url string, host string, maxDepth int) []string {
	visited := make(map[string]struct{}, 300)
	var queue map[string]struct{}
	nextQueue := map[string]struct{}{
		url: {},
	}

	for i := 1; i <= maxDepth; i++ {
		queue, nextQueue = nextQueue, make(map[string]struct{}, 300)
		if len(queue) == 0 {
			break
		}
		println("\n\nDepth: ", i, " Queue: ", len(queue))
		for url := range queue {
			if _, ok := visited[url]; ok {
				println("Already visited")
				continue
			}
			visited[url] = struct{}{}
			for _, link := range visit(url, host) {
				nextQueue[link] = struct{}{}
			}

		}
	}
	ret := make([]string, 0, len(visited))
	for url := range visited {
		ret = append(ret, url)
	}
	return ret
}

func visit(target string, host string) []string {
	println("Visiting: " + target)
	res, err := http.Get(target)
	if err != nil {
		return []string{}
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		println(res.StatusCode)

		return []string{}
	}
	links, err := lp.ParseLinks(res.Body)
	if err != nil {
		println(err)
		return []string{}
	}
	print("Found ", len(links))
	return filterLinks(links, host)

}

func filterLinks(links []lp.Link, host string) []string {
	filteredLinks := []string{}

	for _, v := range links {
		if string(v.Href[0]) == "/" {
			link := strings.TrimSuffix(host+v.Href, "/")
			filteredLinks = append(filteredLinks, link)
			continue
		}

		if strings.HasPrefix(v.Href, host) {

			filteredLinks = append(filteredLinks, strings.TrimSuffix(v.Href, "/"))
		}
	}
	println(" After filter ", len(filteredLinks))
	return filteredLinks
}
