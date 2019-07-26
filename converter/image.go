package converter

import (
	"fmt"
	"image"
	"net/http"
	"strings"

	"golang.org/x/net/html"
)

// ImageConverter is convert the image html.node to amp html.node return true if success
func ImageConverter(n *html.Node) bool {
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
		widthAttr = GetAttributeByName("width", n)

		n.Attr = append(n.Attr, html.Attribute{"", "height", ""})
		heightAttr = GetAttributeByName("height", n)
	}

	if widthAttr.Val == "" || strings.Contains(widthAttr.Val, "%") {
		layoutResponsive = true
	}

	if heightAttr.Val == "" || strings.Contains(heightAttr.Val, "%") {
		layoutResponsive = true
	}

	if layoutResponsive {
		image := getImage(attr.Val)
		if image == nil {
			return false
		}
		bounds := image.Bounds()
		size := bounds.Size()

		if widthAttr.Val != "" && strings.Contains(widthAttr.Val, "%") {
			// calculate ratio
		} else {
			widthAttr.Val = fmt.Sprintf("%d", size.X)
			heightAttr.Val = fmt.Sprintf("%d", size.Y)
		}

	}

	return true
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
