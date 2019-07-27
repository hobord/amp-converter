package converter

import (
	"bytes"
	"log"
	"reflect"
	"strings"
	"sync"

	"golang.org/x/net/html"

	cache "github.com/hobord/amp-converter/cache"
)

// AllowedElements contains all alllowed node elements
var AllowedElements = []string{
	"A",
	"ABBR",
	"ACRONYM",
	"ADDRESS",
	"B",
	"BDO",
	"BIG",
	"BLOCKQUOTE",
	"BODY",
	"BR",
	"BUTTON",
	"CANVAS",
	"CAPTION",
	"CENTER",
	"CITE",
	"CODE",
	"COL",
	"COLGROUP",
	"DD",
	"DEL",
	"DFN",
	"DIR",
	"DIV",
	"DL",
	"DT",
	"EM",
	"FIELDSET",
	"FONT",
	"FORM",
	"H1",
	"H2",
	"H3",
	"H4",
	"H5",
	"H6",
	"HR",
	"HTML",
	"I",
	"IMG",
	"INPUT",
	"INS",
	"ISINDEX",
	"KBD",
	"LABEL",
	"LEGEND",
	"LI",
	"LINK",
	"MAP",
	"MENU",
	"NOSCRIPT",
	"OL",
	"OPTGROUP",
	"OPTION",
	"P",
	"PRE",
	"Q",
	"S",
	"SAMP",
	"SELECT",
	"SMALL",
	"SPAN",
	"STRIKE",
	"STRONG",
	"SUB",
	"SUP",
	"TABLE",
	"TBODY",
	"TD",
	"TEXTAREA",
	"TFOOT",
	"TH",
	"THEAD",
	"TR",
	"TT",
	"U",
	"UL",
	"VAR"}

type toConvert struct {
	node        *html.Node
	baseUrl     string
	ch          cache.Cache
	deleteNodes *[]*html.Node
	wg          *sync.WaitGroup
}

// Converter is convert html to amp
func Converter(htmlDocument string, baseUrl string, ch cache.Cache) string {
	doc, err := html.Parse(strings.NewReader(htmlDocument))
	if err != nil {
		log.Fatal(err) // TODO: do not use fatal
	}

	var wg sync.WaitGroup
	deleteNodes := []*html.Node{}
	// toHeader := []strings

	// convertNode(doc, baseUrl, ch, &deleteNodes, &wg)
	convertNode(toConvert{
		doc,
		baseUrl,
		ch,
		&deleteNodes,
		&wg})
	wg.Wait()

	for _, n := range deleteNodes {
		removeNode(n)
	}

	buf := new(bytes.Buffer)
	html.Render(buf, doc)
	return buf.String()
}

// Converter convert the html.Node tree to AMP node tree.
// func convertNode(n *html.Node, baseUrl string, ch cache.Cache, deleteNodes *[]*html.Node, wg *sync.WaitGroup) {
func convertNode(convertParams toConvert) {
	switch convertParams.node.Type {
	case html.ErrorNode:
		*convertParams.deleteNodes = append(*convertParams.deleteNodes, convertParams.node)
	// case html.TextNode:
	// case html.DocumentNode:
	case html.ElementNode:
		// check attributes
		attributes := []html.Attribute{}
		for i := range convertParams.node.Attr {
			// javascript is not allowed
			if strings.HasPrefix(strings.ToLower(convertParams.node.Attr[i].Val), "javascript:") {
				continue
			}

			// Attribute names starting with on (such as onclick or onmouseover) are disallowed in AMP HTML.
			if strings.HasPrefix(strings.ToLower(convertParams.node.Attr[i].Key), "on") {
				continue
			}

			// XML-related attributes, such as xmlns, xml:lang, xml:base, and xml:space are disallowed in AMP HTML.
			if strings.HasPrefix(strings.ToLower(convertParams.node.Attr[i].Key), "xml") {
				continue
			}
			// Internal AMP attributes prefixed with i-amp- are disallowed in AMP HTML.
			if strings.HasPrefix(strings.ToLower(convertParams.node.Attr[i].Key), "i-amp-") {
				continue
			}

			// check the classess
			if strings.ToLower(convertParams.node.Attr[i].Key) == "class" {
				classes := []string{}
				for _, className := range strings.Fields(convertParams.node.Attr[i].Val) {
					// Internal AMP class names prefixed with -amp- and i-amp- are disallowed in AMP HTML.
					if strings.HasPrefix(className, "-amp-") || strings.HasPrefix(className, "i-amp-") {
						continue
					}
					classes = append(classes, className)
				}
				convertParams.node.Attr[i].Val = strings.Join(classes, " ")
			}

			// Internal AMP IDs prefixed with -amp- and i-amp- are disallowed in AMP HTML.
			if strings.ToLower(convertParams.node.Attr[i].Key) == "id" {
				if strings.HasPrefix(strings.ToLower(convertParams.node.Attr[i].Val), "-amp-") || strings.HasPrefix(strings.ToLower(convertParams.node.Attr[i].Val), "i-amp-") {
					continue
				}
			}

			attributes = append(attributes, convertParams.node.Attr[i])
		}
		convertParams.node.Attr = attributes

		nodeName := strings.ToUpper(convertParams.node.Data)
		// if the node type not allowed then convert or remove
		exists := inArray(nodeName, AllowedElements)
		if exists < 0 {
			switch nodeName {
			case "SCRIPT":
				// script allowed inf the type is  "application/ld+json" or "text/plain"
				allowedTypes := []string{"APPLICATION/LD+JSON", "TEXT/PLAIN"}
				attribute := GetAttributeByName("type", convertParams.node)
				if attribute == nil || inArray(strings.ToUpper(attribute.Val), allowedTypes) < 0 {
					*convertParams.deleteNodes = append(*convertParams.deleteNodes, convertParams.node)
					return
				}
				return
			case "IFRAME":
				// check it is youtube video?
				attribute := GetAttributeByName("src", convertParams.node)
				if attribute != nil || strings.HasPrefix(strings.ToLower(attribute.Val), "https://www.youtube.com/embed") {
					// convert to youtube amp
					if !YoutubeConverter(convertParams.node) {
						// conversion was fail, remove the image
						*convertParams.deleteNodes = append(*convertParams.deleteNodes, convertParams.node)
					}
				}

			}
			*convertParams.deleteNodes = append(*convertParams.deleteNodes, convertParams.node)
			return
		}
		// Some node type is partial allowed need conversion or check
		switch nodeName {
		case "INPUT":
			// "<input[type=image]>, <input[type=button]>, <input[type=password]>, <input[type=file]>" are invalid
			disallowedTypes := []string{"IMAGE", "BUTTON", "PASSWORD", "FILE"}
			attribute := GetAttributeByName("type", convertParams.node)
			if attribute != nil && inArray(strings.ToUpper(attribute.Val), disallowedTypes) >= 0 {
				*convertParams.deleteNodes = append(*convertParams.deleteNodes, convertParams.node)
				return
			}
		case "IMG":
			// convert image tag to amp-img
			convertParams.wg.Add(1)
			go func(n *html.Node, baseUrl string, ch cache.Cache) {
				if !ImageConverter(n, baseUrl, ch) {
					// image conversion was fail, remove the image
					*convertParams.deleteNodes = append(*convertParams.deleteNodes, convertParams.node)
				}
				convertParams.wg.Done()
			}(convertParams.node, convertParams.baseUrl, convertParams.ch)
			return
		}

		// case html.CommentNode:
		// case html.DoctypeNode:
		// case html.scopeMarkerNode:
	}

	for c := convertParams.node.FirstChild; c != nil; c = c.NextSibling {
		convertNode(toConvert{
			c,
			convertParams.baseUrl,
			convertParams.ch,
			convertParams.deleteNodes,
			convertParams.wg})
	}
}

func inArray(val interface{}, array interface{}) int {

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				return i
			}
		}
	}

	return -1
}

// GetAttributeByName return the node attribute
func GetAttributeByName(look string, n *html.Node) *html.Attribute {

	for i := range n.Attr {
		if strings.ToUpper(n.Attr[i].Key) == strings.ToUpper(look) {
			return &n.Attr[i]
		}
	}

	return nil
}

func removeNode(n *html.Node) {
	par := n.Parent
	if par != nil {
		par.RemoveChild(n)
	} else {
		panic("\nNode to remove has no Parent\n") // TODO: do not panic
	}
}
