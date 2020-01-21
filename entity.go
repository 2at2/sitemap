package sitemap

import "encoding/xml"

type Index struct {
	XMLName xml.Name `xml:"sitemapindex"`
	Sitemap []struct {
		Loc     string `xml:"loc"`
		LastMod string `xml:"lastmod"`
	} `xml:"sitemap"`
}

type UrlSet struct {
	XMLName xml.Name `xml:"urlset"`
	URL     []struct {
		Loc        string  `xml:"loc"`
		LastMod    string  `xml:"lastmod"`
		ChangeFreq string  `xml:"changefreq"`
		Priority   float32 `xml:"priority"`
	} `xml:"url"`
}
