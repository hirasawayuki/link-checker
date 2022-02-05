package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/hirasawayuki/link-checker/httprequest"
)

var pageURL string

func init() {
	flag.StringVar(&pageURL, "u", "", "Check page URL.")
}

func main() {
	flag.Parse()
	if pageURL == "" {
		fmt.Println("URL is required.")
		return
	}

	checkResult, err := httprequest.CheckPage(pageURL)
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Link")
	for _, r := range checkResult.AnchorResults {
		fmt.Println(r)
	}

	fmt.Println("Image")
	for _, r := range checkResult.ImgResults {
		fmt.Println(r)
	}
}
