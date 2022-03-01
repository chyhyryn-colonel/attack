package main

import (
	"crypto/md5"
	"encoding/hex"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
	"fmt"
)

const fetchURL = "https://raw.githubusercontent.com/chyhyryn-colonel/attack/main/urls"

func FetchURLs() []string {
	resp, err := http.Get(fetchURL)
	if err != nil {
		fmt.Printf("can't fetch URL list: %v", err)
		return nil
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("can't fetch URL list: %v", err)
		return nil
	}
	urls := strings.Split(string(body), "\n")
	sort.Strings(urls)
	return urls[1:]
}

func Checksum(urls []string) string {
	hash := md5.New()
	for _, u := range urls {
		hash.Write([]byte(u))
	}
	return hex.EncodeToString(hash.Sum(nil))
}
