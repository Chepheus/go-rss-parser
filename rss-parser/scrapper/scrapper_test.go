package scrapper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestVisitRss(t *testing.T) {
	srv := getHttpTestServer()
	srv.Start()
	defer srv.Close()

	rssScrapper := NewRssScrapper()
	posts, err := rssScrapper.VisitRss(srv.URL + "/xml")
	if err != nil {
		t.Error(err.Error())
	}

	if len(posts) == 0 {
		t.Error("posts are empty")
	}

	if posts[0].Title != "Test title" {
		t.Error("post title is wrong")
	}

	if posts[0].Description != "Test Description." {
		t.Error("post description is wrong")
	}

	if posts[0].Link != "https://www.bbc.com/news/articles/test_link" {
		t.Error("post link is wrong")
	}

	if posts[0].Thumbnail != "https://test.com/news/test.jpg" {
		t.Error("post thumbnail is wrong")
	}

	if posts[0].PubDate != "Thu, 20 Jun 2024 15:57:36 GMT" {
		t.Error("post pub_date is wrong")
	}
}

func getHttpTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/xml", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/xml")
		w.Write([]byte(`<rss xmlns:dc="http://purl.org/dc/elements/1.1/" xmlns:content="http://purl.org/rss/1.0/modules/content/" xmlns:atom="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" version="2.0">
<channel>
	<item>
		<title>
			<![CDATA[ Test title ]]>
		</title>
		<description>
			<![CDATA[ Test Description. ]]>
		</description>
		<link>https://www.bbc.com/news/articles/test_link</link>
		<pubDate>Thu, 20 Jun 2024 15:57:36 GMT</pubDate>
		<media:thumbnail width="240" height="135" url="https://test.com/news/test.jpg"/>
	</item>
</channel>
</rss>
		`))
	})

	return httptest.NewUnstartedServer(mux)
}
