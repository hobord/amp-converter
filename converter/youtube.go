package converter

import (
	"fmt"
	"math"
	"regexp"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const defaultAspectRatio = 0.5625

func YoutubeConverter(n *html.Node, ampComponents *AmpComponents) bool {
	attr := GetAttributeByName("src", n)
	if attr == nil {
		return false
	}
	if attr.Val == "" {
		return false
	}
	videoID := findVideoID(attr.Val)
	if videoID == "" {
		return false
	}

	layoutResponsive := false
	widthAttr := GetAttributeByName("width", n)
	if widthAttr == nil {
		layoutResponsive = true
		AddAttribute(n, "width", "")
		widthAttr = GetAttributeByName("width", n)

	}
	if widthAttr.Val == "" || strings.Contains(widthAttr.Val, "%") {
		widthAttr.Val = "560"
		layoutResponsive = true
	}

	heightAttr := GetAttributeByName("height", n)
	if heightAttr == nil {
		AddAttribute(n, "height", "")
		heightAttr = GetAttributeByName("height", n)
	}
	if heightAttr.Val == "" || strings.Contains(heightAttr.Val, "%") {
		width, err := strconv.ParseFloat(widthAttr.Val, 64)
		if err != nil {
			heightAttr.Val = "315"
		} else {
			heightAttr.Val = fmt.Sprintf("%d", int(math.Round((defaultAspectRatio * width))))
		}
	}

	if layoutResponsive {
		AddAttribute(n, "layout", "responsive")
	}
	n.Data = "amp-youtube"
	ampComponents.Add(YoutubeAmpComponent)
	return true
}

func findVideoID(url string) string {
	videoID := url
	if strings.Contains(videoID, "youtu") || strings.ContainsAny(videoID, "\"?&/<%=") {
		reList := []*regexp.Regexp{
			regexp.MustCompile(`(?:v|embed|watch\?v)(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`(?:=|/)([^"&?/=%]{11})`),
			regexp.MustCompile(`([^"&?/=%]{11})`),
		}
		for _, re := range reList {
			if isMatch := re.MatchString(videoID); isMatch {
				subs := re.FindStringSubmatch(videoID)
				videoID = subs[1]
			}
		}
	}
	if strings.ContainsAny(videoID, "?&/<%=") {
		return ""
	}
	if len(videoID) < 10 {
		return ""
	}
	return videoID
}
