package linkchecker

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
)

func Exec(pageURL string, all bool, interval int) error {
	fmt.Printf("Check Page URL: %s\n", pageURL)

	iostream := NewIOStream()
	iostream.StartIndicator()
	defer iostream.StopIndicator()

	u, err := url.Parse(pageURL)
	if err != nil {
		return fmt.Errorf("[ERROR] Parse URL failed. Plese check page url. (url=%s)", pageURL)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return fmt.Errorf("[ERROR] Request failed. err=%s", err)
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("[ERROR] Request failed. url=%s, HTTP Status=%d", u, resp.StatusCode)
	}

	body := resp.Body
	defer resp.Body.Close()
	checkResult, err := CheckPage(body, u.Host, u.Scheme, interval)
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

	font := NewFont()
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

	return nil
}
