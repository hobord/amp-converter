package converter

import (
	// "reflect"
	"bytes"
	"fmt"
	"log"
	"strings"
	"testing"
	"time"

	cache "github.com/hobord/amp-converter/cache"
	"golang.org/x/net/html"
)

func TestImageConverter(t *testing.T) {
	var ch cache.Cache
	ch = cache.NewInMemoryCache(0 * time.Second)
	baseURL := "http://base.com/"
	tests := []struct {
		source string
		want   string
	}{
		{
			source: `<image src="http://p.agnihotry.com/images/avatar-icon.jpeg">`,
			want:   `<amp-img src="http://p.agnihotry.com/images/avatar-icon.jpeg" width="542" height="552" layout="responsive"></amp-img>`,
		},
	}
	for _, tt := range tests {
		doc, err := html.Parse(strings.NewReader(tt.source))
		if err != nil {
			log.Fatal(err)
		}

		if !ImageConverter(doc, baseURL, ch) {
			t.Errorf("ImageConverter got unexpected error")
		}
		buf := new(bytes.Buffer)
		// doc.FirstChild.FirstChild.Data = "span"
		html.Render(buf, doc)
		fmt.Println(buf.String())
	}
}
