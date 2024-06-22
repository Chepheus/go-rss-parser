package scrapper

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

var post = WebPost{
	MainImageSrc: "https://test.com/image.jpg",
	TextBlocks: []string{
		"text block 1",
		"text block 2",
	},
}

func TestVisit(t *testing.T) {
	srv := getHttpTestServer()
	srv.Start()
	defer srv.Close()

	rssScrapper := NewWebScrapper()
	post.ExternalPostLink = srv.URL + "/html"
	p, err := rssScrapper.Visit(post.ExternalPostLink)
	if err != nil {
		t.Error(err.Error())
	}

	if p.ExternalPostLink != post.ExternalPostLink {
		t.Errorf("external_post_link is different: %s not equal %s", post.ExternalPostLink, p.ExternalPostLink)
	}

	if p.MainImageSrc != post.MainImageSrc {
		t.Errorf("main_image_src is different: %s not equal %s", post.MainImageSrc, p.MainImageSrc)
	}

	for i, b := range p.TextBlocks {
		if post.TextBlocks[i] != b {
			t.Errorf("text_block is different: %s not equal %s", post.TextBlocks[i], b)
		}
	}
}

func getHttpTestServer() *httptest.Server {
	mux := http.NewServeMux()

	mux.HandleFunc("/html", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`
<html>
	<body>
		<div id="main-content">
			<article>
				<img src="https://test.com/image.jpg" loading="eager">
				<div data-component="text-block">
					<p>text block 1</p>
					<p>text block 2</p>
				</div>
			</article>
		</div>			
	</body>
</html>
		`))
	})

	return httptest.NewUnstartedServer(mux)
}
