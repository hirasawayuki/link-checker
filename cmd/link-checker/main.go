package main

import (
	"flag"
	"log"
	"os"

	linkchecker "github.com/hirasawayuki/link-checker"
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

	err := linkchecker.Exec(pageURL, all, interval)
	if err != nil {
		log.Fatalln(err)
	}
	os.Exit(0)
}
