package converter

import "golang.org/x/net/html"

// IframeConverter convert html iframe to amp-iframe
func IframeConverter(n *html.Node, ampComponents *AmpComponents) bool {
	attr := GetAttributeByName("src", n)
	if attr == nil {
		return false
	}
	if attr.Val == "" {
		return false
	}

	n.Data = "amp-iframe"
	ampComponents.Add(IframeAmpComponent)
	return true
}
