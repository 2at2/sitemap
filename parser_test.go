package sitemap

import (
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestParserSimple(t *testing.T) {
	data := []byte(`<?xml version="1.0" encoding="UTF-8"?>
<!--generator='jetpack-7.8.1'-->
<?xml-stylesheet type="text/xsl" href="//example.com/sitemap.xsl"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">
    <url>
        <loc>%addr%/contributors/</loc>
        <lastmod>2019-12-30T09:57:58Z</lastmod>
    </url>
    <url>
        <loc>%addr%/contact/</loc>
        <lastmod>2019-12-30T09:56:13Z</lastmod>
    </url>
    <url>
        <loc>%addr%/form/</loc>
        <lastmod>2019-11-11T23:09:44Z</lastmod>
    </url>
    <url>
        <loc>%addr%/form/</loc>
        <lastmod>2019-11-11T23:09:44Z</lastmod>
    </url>
</urlset>`)

	parser := parser{client: http.DefaultClient}
	urls, err := parser.parse(data, 0)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 3, len(urls))
}

type SitemapHandler struct {
}

func (m SitemapHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	var body string

	switch path {
	case "/":
		body = `<?xml version="1.0" encoding="UTF-8"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <sitemap>
        <loc>%addr%/sitemap-1.xml</loc>
        <lastmod>2020-01-18T03:48:22Z</lastmod>
    </sitemap>
    <sitemap>
        <loc>%addr%/sitemap-2.xml</loc>
        <lastmod>2020-01-17T12:01:52Z</lastmod>
    </sitemap>
</sitemapindex>`

	case "/sitemap-1.xml":
		body = `<?xml version="1.0" encoding="UTF-8"?>
<!--generator='jetpack-7.8.1'-->
<?xml-stylesheet type="text/xsl" href="//example.com/sitemap.xsl"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">
    <url>
        <loc>%addr%/contributors/</loc>
        <lastmod>2019-12-30T09:57:58Z</lastmod>
    </url>
    <url>
        <loc>%addr%/contact/</loc>
        <lastmod>2019-12-30T09:56:13Z</lastmod>
    </url>
    <url>
        <loc>%addr%/form/</loc>
        <lastmod>2019-11-11T23:09:44Z</lastmod>
    </url>
    <url>
        <loc>%addr%/form/</loc>
        <lastmod>2019-11-11T23:09:44Z</lastmod>
    </url>
</urlset>`

	case "/sitemap-2.xml":
		body = `<?xml version="1.0" encoding="UTF-8"?>
<!--generator='jetpack-7.8.1'-->
<?xml-stylesheet type="text/xsl" href="//detroitisit.com/sitemap-index.xsl"?>
<sitemapindex xmlns="http://www.sitemaps.org/schemas/sitemap/0.9">
    <sitemap>
        <loc>%addr%/sitemap-3.xml</loc>
        <lastmod>2019-12-24T17:11:27Z</lastmod>
    </sitemap>
</sitemapindex>`

	case "/sitemap-3.xml":
		body = `<?xml version="1.0" encoding="UTF-8"?>
<!--generator='jetpack-7.8.1'-->
<?xml-stylesheet type="text/xsl" href="//example.com/sitemap.xsl"?>
<urlset xmlns="http://www.sitemaps.org/schemas/sitemap/0.9" xmlns:xsi="http://www.w3.org/2001/XMLSchema-instance"
        xsi:schemaLocation="http://www.sitemaps.org/schemas/sitemap/0.9 http://www.sitemaps.org/schemas/sitemap/0.9/sitemap.xsd">
    <url>
        <loc>%addr%/news/1</loc>
        <lastmod>2019-12-30T09:57:58Z</lastmod>
    </url>
    <url>
        <loc>%addr%/news/1</loc>
        <lastmod>2019-12-30T09:57:58Z</lastmod>
    </url>
	<url>
        <loc>%addr%/news/2</loc>
        <lastmod>2019-12-30T09:57:58Z</lastmod>
    </url>
</urlset>`

	default:
		w.WriteHeader(404)
		return
	}

	body = strings.ReplaceAll(body, "%addr%", "http://"+r.Host)

	w.Header().Add("Content-Type", "application/xml")
	w.Write([]byte(body))
}

func TestParserHttp(t *testing.T) {
	server := httptest.NewServer(SitemapHandler{})
	defer server.Close()

	urls, err := Get(server.URL)
	if !assert.NoError(t, err) {
		t.FailNow()
	}

	assert.Equal(t, 5, len(urls))
}
