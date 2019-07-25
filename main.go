package main

import (
	"bytes"
	"errors"
	"fmt"
	"log"
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

func main() {
	fmt.Println("App started")
	htmlDocument := `
		<div id="mainDivId">
			<h1>Head 1</h1>
			<p>
				paragraph text wit some <em class="underline">inline</em> elements
			</p>
		</div>
	`
	fmt.Println(htmlDocument)
	doc, err := html.Parse(strings.NewReader(htmlDocument))
	if err != nil {
		log.Fatal(err)
	}

	domConverter(doc)
	buf := new(bytes.Buffer)
	html.Render(buf, doc)
	fmt.Println(buf.String())
}

func domConverter(n *html.Node) {
	switch n.Type {
	// case html.ErrorNode:
	// case html.TextNode:
	// case html.DocumentNode:
	case html.ElementNode:
		/*
			check the node type
				script allowed inf the type is  "application/ld+json" or "text/plain"
				"<input[type=image]>, <input[type=button]>, <input[type=password]>, <input[type=file]>" are invalid
				<A href attribute value must not begin with javascript:
			if it is not in the whitelist then convert it.  If set, the target attribute value must be _blank

			check the attributes
				Attribute names starting with on (such as onclick or onmouseover) are disallowed in AMP HTML.
				XML-related attributes, such as xmlns, xml:lang, xml:base, and xml:space are disallowed in AMP HTML.
				Internal AMP attributes prefixed with i-amp- are disallowed in AMP HTML.

			check classes if present
				Internal AMP class names prefixed with -amp- and i-amp- are disallowed in AMP HTML.
			IDs
				Internal AMP IDs prefixed with -amp- and i-amp- are disallowed in AMP HTML.
			Links
				The javascript: schema is disallowed.

		*/
		nodeName := strings.ToUpper(n.Data)
		exists, _ := in_array(nodeName, AllowedElements)
		if !exists {
			switch nodeName {
			case "SCRIPT":
				attribute, error := getAttributeByName("type", n)
				if error != nil || (attribute.Val != "application/ld+json") {
					panic("implement")
				}
			}
		}
		// case html.CommentNode:
		// case html.DoctypeNode:
		// case html.scopeMarkerNode:

	}

	fmt.Println(n.Type)
	fmt.Println(n.Data)
	fmt.Println(n.Attr)
	if n.Data == "Head 1" {
		n.Data = "Head modifyed"
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		domConverter(c)
	}
}

func in_array(val interface{}, array interface{}) (exists bool, index int) {
	exists = false
	index = -1

	switch reflect.TypeOf(array).Kind() {
	case reflect.Slice:
		s := reflect.ValueOf(array)

		for i := 0; i < s.Len(); i++ {
			if reflect.DeepEqual(val, s.Index(i).Interface()) == true {
				index = i
				exists = true
				return
			}
		}
	}

	return
}

func getAttributeByName(look string, n *html.Node) (html.Attribute, error) {
	attributes := n.Attr
	for _, attr := range attributes {
		if strings.ToUpper(attr.Key) == strings.ToUpper(look) {
			return attr, nil
		}
	}

	return html.Attribute{}, errors.New("Not found")
}
