package sitemap

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
)

type Parser interface {
	Get(url string) ([]string, error)
}

var DefaultParser = NewParser(nil, 100)

// parser main parser struct
type parser struct {
	client   *http.Client
	maxDepth int
}

// NewParser returns new parser struct
func NewParser(client *http.Client, maxDepth int) Parser {
	if client == nil {
		client = http.DefaultClient
	}

	return &parser{client: client, maxDepth: maxDepth}
}

// Get returns sitemap urls
func (p *parser) Get(url string) ([]string, error) {
	return p.get(url, 0)
}

func Get(url string) ([]string, error) {
	return DefaultParser.Get(url)
}

// get returns sitemap urls
func (p *parser) get(url string, depth int) ([]string, error) {
	// Throttle
	if p.maxDepth > 0 && depth > p.maxDepth {
		return nil, nil
	}

	resp, err := p.client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	bts, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return p.parse(bts, depth)
}

func (p *parser) parse(data []byte, depth int) ([]string, error) {
	var urls []string

	if set, err := p.parseUrlSet(data); err == nil {
		for _, url := range set.URL {
			urls = append(urls, url.Loc)
		}
	} else if index, err := p.parseIndex(data); err == nil {
		for _, m := range index.Sitemap {
			list, err := p.get(m.Loc, depth+1)
			if err != nil {
				return nil, err
			}
			urls = append(urls, list...)
		}
	} else {
		return nil, err
	}

	return p.sliceToUnique(urls), nil
}

// parseUrlSet returns sitemap UrlSet struct
func (p *parser) parseUrlSet(data []byte) (set UrlSet, err error) {
	err = xml.Unmarshal(data, &set)
	return
}

// parseUrlSet returns sitemap Index struct
func (p *parser) parseIndex(data []byte) (index Index, err error) {
	err = xml.Unmarshal(data, &index)
	return
}

// sliceToUnique returns unique slice
func (p *parser) sliceToUnique(source []string) []string {
	if len(source) == 0 {
		return source
	}

	var result []string
	m := make(map[string]string)
	for _, x := range source {
		m[x] = x
	}
	for _, x := range m {
		result = append(result, x)
	}

	return result
}
