package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/briandowns/spinner"
	"github.com/hirasawayuki/link-checker/httprequest"
)

var pageURL string
var all bool

func init() {
	flag.StringVar(&pageURL, "u", "", "Check page URL.")
	flag.BoolVar(&all, "a", false, "Display all status.")
}

func main() {
	s := spinner.New(spinner.CharSets[11], 100*time.Millisecond)
	s.Start()

	flag.Parse()
	if pageURL == "" {
		log.Fatalln("URL is required.")
		return
	}

	checkResult, err := httprequest.CheckPage(pageURL)
	s.Stop()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Link")
	for _, r := range checkResult.AnchorResults {
		if !all && r.Status == http.StatusOK {
			continue
		}
		fmt.Println(r.Text)
	}

	fmt.Println("Image")
	for _, r := range checkResult.ImgResults {
		if !all && r.Status == http.StatusOK {
			continue
		}
		fmt.Println(r.Text)
	}
}
