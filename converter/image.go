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
func ImageConverter(n *html.Node, baseURL string, ch cache.Cache) bool {
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

	srcAttr := GetAttributeByName("src", n)
	if srcAttr == nil {
		return false
	}
	if srcAttr.Val == "" {
		return false
	}
	if !strings.HasPrefix(strings.ToLower(srcAttr.Val), "http") {
		srcAttr.Val = baseURL + srcAttr.Val
	}

	newAttr := []html.Attribute{}

	// Default the image is not responsive and we try using the attr sizes parameters
	layoutAttr := GetAttributeByName("layout", n)
	if layoutAttr == nil {
		AddAttribute(n, "layout", "")
		layoutAttr = GetAttributeByName("layout", n)
	}
	widthAttr := GetAttributeByName("width", n)
	heightAttr := GetAttributeByName("height", n)

	if layoutAttr.Val != "fill" && layoutAttr.Val != "nodisplay" {
		var size image.Point
		var er error
		// we need the height
		if heightAttr == nil || heightAttr.Val == "" || strings.Contains(heightAttr.Val, "%") {
			size, er = getImageSize(srcAttr.Val, ch)
			if er != nil {
				return false
			}
			AddAttribute(n, "height", fmt.Sprintf("%d", size.Y))
			heightAttr = GetAttributeByName("height", n)

			if layoutAttr.Val == "" {
				layoutAttr.Val = "responsive"
			}
		}

		// "fixed-height" is not required the width
		if layoutAttr.Val == "" || layoutAttr.Val == "fixed" || layoutAttr.Val == "intrinsic" || layoutAttr.Val == "responsive" {
			// width and height are mandatory
			if widthAttr == nil || widthAttr.Val == "" || strings.Contains(widthAttr.Val, "%") {
				if size.X == 0 {
					size, er = getImageSize(srcAttr.Val, ch)
					if er != nil {
						return false
					}
				}
				AddAttribute(n, "width", fmt.Sprintf("%d", size.X))
				widthAttr = GetAttributeByName("width", n)
				if layoutAttr.Val == "" {
					layoutAttr.Val = "responsive"
				}
			}
		}
	}

	newAttr = append(newAttr, *srcAttr)
	if layoutAttr.Val != "fixed-height" && layoutAttr.Val != "fill" && layoutAttr.Val != "nodisplay" {
		if widthAttr != nil {
			newAttr = append(newAttr, *widthAttr)
		}
	}
	if heightAttr != nil && layoutAttr.Val != "fill" && layoutAttr.Val != "nodisplay" {
		newAttr = append(newAttr, *heightAttr)
	}
	if layoutAttr != nil && layoutAttr.Val != "" {
		newAttr = append(newAttr, *layoutAttr)
	}
	titleAttr := GetAttributeByName("title", n)
	if titleAttr != nil && titleAttr.Val != "" {
		newAttr = append(newAttr, *titleAttr)
	}
	altAttr := GetAttributeByName("alt", n)
	if altAttr != nil && altAttr.Val != "" {
		newAttr = append(newAttr, *altAttr)
	}
	classAttr := GetAttributeByName("class", n)
	if classAttr != nil && classAttr.Val != "" {
		newAttr = append(newAttr, *classAttr)
	}
	mediaAttr := GetAttributeByName("media", n)
	if mediaAttr != nil && mediaAttr.Val != "" {
		newAttr = append(newAttr, *mediaAttr)
	}
	sizesAttr := GetAttributeByName("sizes", n)
	if sizesAttr != nil && sizesAttr.Val != "" {
		newAttr = append(newAttr, *sizesAttr)
	}
	srcsetAttr := GetAttributeByName("srcset", n)
	if srcsetAttr != nil && srcsetAttr.Val != "" {
		newAttr = append(newAttr, *srcsetAttr)
	}
	heightsAttr := GetAttributeByName("heights", n)
	if heightsAttr != nil && heightsAttr.Val != "" {
		newAttr = append(newAttr, *heightsAttr)
	}
	noloadingAttr := GetAttributeByName("noloading", n)
	if noloadingAttr != nil {
		newAttr = append(newAttr, *noloadingAttr)
	}

	n.Attr = newAttr
	n.Data = "amp-img"
	return true
}

func getImageSize(url string, ch cache.Cache) (image.Point, error) {
	size := image.Point{0, 0}

	h := sha256.New()
	h.Write([]byte(url))
	key := fmt.Sprintf("%x", h.Sum(nil))
	err := errors.New("No Cahce")
	if ch != nil {
		err = ch.Get(key, &size)
	}
	if err != nil {
		img := getImage(url)
		if img == nil {
			return image.Point{0, 0}, errors.New("Cant get image")
		}
		bounds := img.Bounds()
		size = bounds.Size()
		if ch != nil {
			ch.Set(key, size, 24*time.Hour)
		}
	}
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
