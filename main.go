package main

import (
	"bytes"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/html"

)

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
	case html.ErrorNode:
	case html.TextNode:
	case html.DocumentNode:
	case html.ElementNode:
		/*
			 	check the node type
				if it is not in the whitelist then convert it.
				check the attributes

		*/

	case html.CommentNode:
	case html.DoctypeNode:
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
