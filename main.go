package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hirasawayuki/link-checker/httprequest"
)

var pageURL string

func init() {
	flag.StringVar(&pageURL, "u", "", "Check page URL.")
}

func main() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	flag.Parse()
	if pageURL == "" {
		fmt.Println("URL is required.")
		return
	}

	checkResult, err := httprequest.CheckPage(pageURL)
	s.Stop()
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
