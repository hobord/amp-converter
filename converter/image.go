package converter

import (
	"crypto/sha256"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"net/http"
	"strings"
	"time"

	"github.com/pkg/errors"
	"golang.org/x/net/html"

	cache "github.com/hobord/amp-converter/cache"
)

// ImageConverter is convert the image html.node to amp html.node return true if success
func ImageConverter(n *html.Node, baseUrl string, ch cache.Cache) bool {
	/**
	 	Check the image sizes
		If width is not set or % type then the image layout is responsive and You should download for get the sizes
		If it is set and not % then the image is not responsive, use the image sizes
		use the srcset attribute if is exists
		use the media attribute if is exists
		use the sizes attribute if is exists
		use the alt attribute if is exists
		use the title attribute if is exists
	*/

	attr := GetAttributeByName("src", n)
	if attr == nil {
		return false
	}
	if attr.Val == "" {
		return false
	}

	layoutResponsive := false
	widthAttr := GetAttributeByName("width", n)
	heightAttr := GetAttributeByName("height", n)
	if widthAttr == nil || heightAttr == nil {
		layoutResponsive = true

		n.Attr = append(n.Attr, html.Attribute{"", "width", ""})
		n.Attr = append(n.Attr, html.Attribute{"", "height", ""})
		widthAttr = GetAttributeByName("width", n)
		heightAttr = GetAttributeByName("height", n)
	}

	if widthAttr.Val == "" || strings.Contains(widthAttr.Val, "%") {
		layoutResponsive = true
	}

	if heightAttr.Val == "" || strings.Contains(heightAttr.Val, "%") {
		layoutResponsive = true
	}

	// ampImg := &html.Node{
	// 	Parent:      n.Parent,
	// 	PrevSibling: n.PrevSibling,
	// 	NextSibling: n.NextSibling,
	// 	Type:        n.Type,
	// 	DataAtom:    n.DataAtom,
	// 	Data:        "amp-img",
	// 	Attr:        []html.Attribute{},
	// }

	if layoutResponsive {
		size, error := getImageSize(attr.Val, ch)
		if error != nil {
			return false
		}

		if widthAttr.Val != "" && strings.Contains(widthAttr.Val, "%") {
			// calculate ratio
		} else {
			widthAttr.Val = fmt.Sprintf("%d", size.X)
			heightAttr.Val = fmt.Sprintf("%d", size.Y)
		}

		n.Attr = append(n.Attr, html.Attribute{"", "layout", "responsive"}) // TODO: check if is already exists
	}

	n.Data = "amp-img"
	return true
}

func getImageSize(url string, ch cache.Cache) (image.Point, error) {
	size := image.Point{0, 0}

	h := sha256.New()
	h.Write([]byte(url))
	key := fmt.Sprintf("%x", h.Sum(nil))

	err := ch.Get(key, &size)
	if err != nil {
		img := getImage(url)
		if img == nil {
			return image.Point{0, 0}, errors.New("Cant get image")
		}
		bounds := img.Bounds()
		size = bounds.Size()
		ch.Set(key, size, 24*time.Hour)
	}
	fmt.Println(size)
	return size, nil
}

func getImage(url string) image.Image {
	res, err := http.Get(url)
	if err != nil || res.StatusCode != 200 {
		return nil
	}
	defer res.Body.Close()
	m, _, err := image.Decode(res.Body)
	if err != nil {
		return nil
	}
	return m
}
