package httprequest

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/hirasawayuki/link-checker/html"
)

const (
	defaultRequestSpeedMillisecond = 100
	maxConnectionCount             = 10
)

type CheckResults struct {
	AnchorResults []*Result
	ImgResults    []*Result

	mux sync.Mutex
}

type Result struct {
	Text   string
	Status int
}

func (cr *CheckResults) append(node html.Node, result *Result) {
	cr.mux.Lock()
	defer cr.mux.Unlock()

	switch v := node.(type) {
	case *html.AnchorNode:
		cr.AnchorResults = append(cr.AnchorResults, result)
	case *html.ImgNode:
		cr.ImgResults = append(cr.ImgResults, result)
	default:
		fmt.Println(v)
	}
}

func CheckPage(pageURL string) (*CheckResults, error) {
	check := &CheckResults{}
	u, err := url.Parse(pageURL)
	if err != nil {
		fmt.Printf("[ERROR] Parse URL failed. Plese check page url. (url=%s)\n", pageURL)
		return nil, fmt.Errorf("[ERROR] Parse URL failed. Plese check page url. (url=%s)\n", pageURL)
	}

	resp, err := http.Get(u.String())
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		return nil, fmt.Errorf("[ERROR] Request failed. url=%s, HTTP Status=%d\n", u, resp.StatusCode)
	}

	ns, err := html.Parse(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Parse HTML failed. (err=%w)\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	semaphore := make(chan struct{}, maxConnectionCount)
	for i := 0; i < maxConnectionCount-1; i++ {
		semaphore <- struct{}{}
	}

	go func() {
		for {
			select {
			case <-time.After(defaultRequestSpeedMillisecond * time.Millisecond):
				<-semaphore
			}
		}
	}()

	errCh := make(chan error, len(ns))
	wg := &sync.WaitGroup{}
	for _, n := range ns {
		url, err := n.URL()
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Unexpected URL. (err=%w)\n", err)
		}
		if url.Host == "" {
			url.Host = u.Host
		}
		if url.Scheme == "" {
			url.Scheme = u.Scheme
		}

		wg.Add(1)
		go func(n html.Node) {
			defer wg.Done()
			semaphore <- struct{}{}
			if err := checkStatus(ctx, url.String(), n, check); err != nil {
				cancel()
				errCh <- err
			}
		}(n)
	}
	wg.Wait()

	close(errCh)
	var errs []error
	for err := range errCh {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		fmt.Printf("\n[WARNING] %s\n", errs[0])
	}

	return check, nil
}

func checkStatus(ctx context.Context, url string, n html.Node, check *CheckResults) error {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	client := &http.Client{
		Timeout: time.Second * 10,
	}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	r := &Result{}
	if resp.StatusCode < http.StatusBadRequest {
		success := color.GreenString("✓")
		r.Text = fmt.Sprintf("%s HTTP Status: %d URL: %s Text(alt): %s", success, resp.StatusCode, url, n)
		r.Status = resp.StatusCode
		check.append(n, r)
		return nil
	}

	failure := color.RedString("X")
	r.Text = fmt.Sprintf("%s HTTP Status: %d URL: %s Text(alt): %s", failure, resp.StatusCode, url, n)
	r.Status = resp.StatusCode
	check.append(n, r)
	return nil
}
