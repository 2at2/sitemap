# Golang sitemap parser

Module parses sitemap url and returns website urls.

## Installation

```
go get github.com/2at2/sitemap
```

## Usage

```
p := parser.NewParser(http.DefaultClient, -1)
urls, err := p.Get("http://example.com/sitemap.xml")

# OR use default parser

urls, err := parser.Get("http://example.com/sitemap.xml")
```
