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

	if !all {
		var i, j int
		for _, r := range checkResult.AnchorResults {
			if r.Status >= http.StatusBadRequest {
				checkResult.AnchorResults[i] = r
				i++
			}
		}
		checkResult.AnchorResults = checkResult.AnchorResults[:i]
		for _, r := range checkResult.ImgResults {
			if r.Status >= http.StatusBadRequest {
				checkResult.ImgResults[j] = r
				j++
			}
		}
		checkResult.ImgResults = checkResult.ImgResults[:j]
	}

	fmt.Println("Link")
	if len(checkResult.AnchorResults) == 0 {
		fmt.Println("All checks have passed.")
	} else {
		for _, r := range checkResult.AnchorResults {
			fmt.Println(r.Text)
		}
	}

	fmt.Printf("\nImage\n")
	if len(checkResult.ImgResults) == 0 {
		fmt.Println("All checks have passed.")
	} else {
		for _, r := range checkResult.ImgResults {
			fmt.Println(r.Text)
		}
	}
}
