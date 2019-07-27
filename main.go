package main

import (
	// "errors"
	"fmt"
	"time"

	cache "github.com/hobord/amp-converter/cache"
	amp "github.com/hobord/amp-converter/converter"

)

func main() {
	fmt.Println("App started")
	htmlDocument := `
		<div id="mainDivId">
			<h1 i-amp-="dsds">Head 1</h1>
			<script>
				<p>dsad</p>
			</script>
			<script type="application/ld+json">
				<p>dsad</p>
			</script>
			<applet>java applet</applet>
			<a href="javascript: alert()">hello</a>
			<input>
			<input type="button">
			<image src="http://p.agnihotry.com/images/avatar-icon.jpeg">
			<span id="i-amp-">dsadas</span>
			<p onmoUse="dss" class="-amp- bold">
			<iframe width="560" height="315" src="https://www.youtube.com/embed/9nqcFhD3wz0" frameborder="0" allow="accelerometer; autoplay; encrypted-media; gyroscope; picture-in-picture" allowfullscreen></iframe>
			<p xml="saa">fd</p>
				<image src="/images/avatar-icon.jpeg">
				paragraph text wit some <em class="underline">inline</em> elements
			</p>
		</div>
	`
	baseURL := "http://p.agnihotry.com"
	var ch cache.Cache
	ch = newInMemoryCache(60 * time.Second)
	fmt.Println(amp.Converter(htmlDocument, baseURL, ch))
}

var newInMemoryCache = func(defaultExpiration time.Duration) cache.Cache {
	return cache.NewInMemoryCache(defaultExpiration)
}
