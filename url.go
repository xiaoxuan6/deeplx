package deeplx

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"
	"time"
)

func fetchUri() string {
	if len(urls) < 1 {

		var wgs sync.WaitGroup
		wgs.Add(2)
		for i, url := range []string{
			"https://github-mirror.us.kg/https://github.com/ycvk/deeplx-local/blob/windows/url.txt",
			"https://github-mirror.us.kg/https://github.com/xiaozhou26/serch_deeplx/blob/main/success.txt",
		} {
			go fetchUrls(&wgs, i, url)
		}
		wgs.Wait()

		for _, url := range targetUrls {
			wg.Add(1)
			go checkUrlVerify(url, &wg)
		}
		wg.Wait()
	}

	randomIndex := randomNum(urls)
	if randomIndex <= len(urls) {
		return urls[randomIndex]
	}

	return urls[0]
}

var client = &http.Client{
	Timeout: 3 * time.Second,
}

func fetchUrls(wg *sync.WaitGroup, k int, url string) {
	defer wg.Done()

	resp, err := client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	r := bufio.NewReader(resp.Body)
	for {
		line, _, errs := r.ReadLine()
		if errs == io.EOF {
			break
		}

		newUrl := string(line)
		if k == 1 {
			newUrl = fmt.Sprintf("%s/translate", newUrl)
		}
		targetUrls = append(targetUrls, newUrl)
	}
}

func checkUrlVerify(url string, wg *sync.WaitGroup) {
	defer wg.Done()

	resp, err := client.Get(strings.ReplaceAll(url, "/translate", ""))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		urls = append(urls, url)
	}
}

func randomNum(slices []string) int {
	urlsLen := len(slices)
	return rand.Intn(urlsLen)
}

var (
	proxyUrls []string
	once      sync.Once
)

func getProxyUrl() string {
	once.Do(func() {
		res, err := client.Get("https://269900.xyz/fetch_http_all")
		if err != nil {
			return
		}

		defer res.Body.Close()
		item, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return
		}

		f := bufio.NewReader(strings.NewReader(string(item)))
		for {
			line, err := f.ReadString('\n')
			if err != nil {
				break
			}

			proxyUrls = append(proxyUrls, line)
		}
	})

	if len(proxyUrls) < 1 {
		return ""
	}

	return proxyUrls[randomNum(proxyUrls)]
}
