package converter

import (
	"reflect"
	"strings"

	"golang.org/x/net/html"
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

// Converter convert the html.Node tree to AMP node tree
func Converter(n *html.Node) *html.Node {
	switch n.Type {
	// case html.ErrorNode:
	// case html.TextNode:
	// case html.DocumentNode:
	case html.ElementNode:
		nodeName := strings.ToUpper(n.Data)
		exists := inArray(nodeName, AllowedElements)
		if exists < 0 {
			switch nodeName {
			case "SCRIPT":
				// script allowed inf the type is  "application/ld+json" or "text/plain"
				allowedTypes := []string{"APPLICATION/LD+JSON", "TEXT/PLAIN"}
				attribute := getAttributeByName("type", n)
				if attribute == nil || inArray(strings.ToUpper(attribute.Val), allowedTypes) < 0 {
					return n
				}
				return nil
			}
			return n
		}

		switch nodeName {
		case "INPUT":
			// "<input[type=image]>, <input[type=button]>, <input[type=password]>, <input[type=file]>" are invalid
			disallowedTypes := []string{"IMAGE", "BUTTON", "PASSWORD", "FILE"}
			attribute := getAttributeByName("type", n)
			if attribute != nil && inArray(strings.ToUpper(attribute.Val), disallowedTypes) >= 0 {
				return n
			}
		case "IMG":
			// convert image tag to amp-img
			// if there is no width and hight then download the image and detect...
			// use go routine with waitgroup parameter will n.
		}

		// check attributes
		attributes := []html.Attribute{}
		for i := range n.Attr {
			// javascript is not allowed
			if strings.HasPrefix(strings.ToLower(n.Attr[i].Val), "javascript:") {
				continue
			}

			// Attribute names starting with on (such as onclick or onmouseover) are disallowed in AMP HTML.
			if strings.HasPrefix(strings.ToLower(n.Attr[i].Key), "on") {
				continue
			}

			// XML-related attributes, such as xmlns, xml:lang, xml:base, and xml:space are disallowed in AMP HTML.
			if strings.HasPrefix(strings.ToLower(n.Attr[i].Key), "xml") {
				continue
			}
			// Internal AMP attributes prefixed with i-amp- are disallowed in AMP HTML.
			if strings.HasPrefix(strings.ToLower(n.Attr[i].Key), "i-amp-") {
				continue
			}

			// check the classess
			if strings.ToLower(n.Attr[i].Key) == "class" {
				classes := []string{}
				for _, className := range strings.Fields(n.Attr[i].Val) {
					// Internal AMP class names prefixed with -amp- and i-amp- are disallowed in AMP HTML.
					if strings.HasPrefix(className, "-amp-") || strings.HasPrefix(className, "i-amp-") {
						continue
					}
					classes = append(classes, className)
				}
				n.Attr[i].Val = strings.Join(classes, " ")
			}

			// Internal AMP IDs prefixed with -amp- and i-amp- are disallowed in AMP HTML.
			if strings.ToLower(n.Attr[i].Key) == "id" {
				if strings.HasPrefix(strings.ToLower(n.Attr[i].Val), "-amp-") || strings.HasPrefix(strings.ToLower(n.Attr[i].Val), "i-amp-") {
					continue
				}
			}

			attributes = append(attributes, n.Attr[i])
		}
		n.Attr = attributes
		// case html.CommentNode:
		// case html.DoctypeNode:
		// case html.scopeMarkerNode:

	}

	// fmt.Println(n.Type)
	// fmt.Println(n.Data)
	// fmt.Println(n.Attr)
	// if n.Data == "Head 1" {
	// 	n.Data = "Head modifyed"
	// }

	var remove *html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if remove != nil {
			removeNode(remove)
			remove = nil
		}
		remove = Converter(c)
	}
	return nil
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

func getAttributeByName(look string, n *html.Node) *html.Attribute {

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
