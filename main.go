package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/hirasawayuki/link-checker/httprequest"
	"github.com/hirasawayuki/link-checker/iostream"
)

var pageURL string
var all bool
var interval int

func init() {
	flag.StringVar(&pageURL, "u", "", "Check page URL.")
	flag.BoolVar(&all, "a", false, "Display all status.")
	flag.IntVar(&interval, "t", 100, "HTTP request interval time. (ms)")
}

func main() {
	flag.Parse()
	if pageURL == "" {
		log.Fatalln("URL is required.")
		return
	}

	fmt.Printf("Check Page URL: %s\n", pageURL)

	iostream := iostream.New()
	iostream.StartIndicator()
	defer iostream.StopIndicator()

	u, err := url.Parse(pageURL)
	if err != nil {
		log.Fatalf("[ERROR] Parse URL failed. Plese check page url. (url=%s)\n", pageURL)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		log.Fatalf("[ERROR] Request failed. err=%s", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Fatalf("[ERROR] Request failed. url=%s, HTTP Status=%d\n", u, resp.StatusCode)
	}

	body := resp.Body
	defer resp.Body.Close()
	checkResult, err := httprequest.CheckPage(body, u.Host, u.Scheme, interval)
	iostream.StopIndicator()

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

	font := iostream.Font()
	fmt.Printf("\n%s\n", font.Bold("[Link]"))
	if len(checkResult.AnchorResults) == 0 {
		fmt.Printf("%s\n", font.Green("✓ All checks have passed."))
	} else {
		for _, r := range checkResult.AnchorResults {
			fmt.Println(r.Text)
		}
	}

	fmt.Printf("\n%s\n", font.Bold("[Image]"))
	if len(checkResult.ImgResults) == 0 {
		fmt.Printf("%s\n", font.Green("✓ All checks have passed."))
	} else {
		for _, r := range checkResult.ImgResults {
			fmt.Println(r.Text)
		}
	}
}
