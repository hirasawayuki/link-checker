package httprequest

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/hirasawayuki/link-checker/html"
)

const (
	maxConnectionCount = 10
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

	switch node.(type) {
	case *html.AnchorNode:
		cr.AnchorResults = append(cr.AnchorResults, result)
	case *html.ImgNode:
		cr.ImgResults = append(cr.ImgResults, result)
	}
}

func CheckPage(r io.Reader, host, scheme string, interval int) (*CheckResults, error) {
	check := &CheckResults{}
	ns, err := html.Parse(r)
	if err != nil {
		return nil, fmt.Errorf("[ERROR] Parse HTML failed. (err=%w)\n", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	semaphore := make(chan struct{}, maxConnectionCount)
	for i := 0; i < maxConnectionCount-1; i++ {
		semaphore <- struct{}{}
	}

	if interval < 50 {
		interval = 50
	}
	go func() {
		for {
			select {
			case <-time.After(time.Duration(interval) * time.Millisecond):
				<-semaphore
			}
		}
	}()

	errCh := make(chan error, len(ns))
	wg := &sync.WaitGroup{}
	for _, n := range ns {
		nodeURL, err := n.URL()
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Unexpected URL. (err=%w)\n", err)
		}
		parsedURL, err := url.Parse(nodeURL)
		if err != nil {
			return nil, fmt.Errorf("[ERROR] Invalid URL. (err=%w)\n", err)
		}
		if parsedURL.Host == "" {
			parsedURL.Host = host
		}
		if parsedURL.Scheme == "" {
			parsedURL.Scheme = scheme
		}

		wg.Add(1)
		go func(n html.Node) {
			defer wg.Done()
			semaphore <- struct{}{}
			if err := checkStatus(ctx, parsedURL.String(), n, check); err != nil {
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
	req, err := http.NewRequest(http.MethodHead, url, nil)
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
