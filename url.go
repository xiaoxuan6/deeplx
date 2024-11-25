package deeplx

import (
	"bufio"
	"fmt"
	"github.com/samber/lo"
	"github.com/xiaoxuan6/deeplx/api/log"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

func LoadBlack(reset bool) {
	b, _ := os.ReadFile("blacklist.txt")
	r := bufio.NewReader(strings.NewReader(string(b)))

	if reset {
		blackList = blackList[:0]
	}

	for {
		line, _, err := r.ReadLine()
		if err == io.EOF {
			break
		}

		newLine := strings.Trim(string(line), "/")
		blackList = append(blackList, newLine)
	}
}

func CheckUrlAndReloadBlack() {
	targetUrls = targetUrls[:0]
	urls = urls[:0]
	_ = fetchUri()

	blackList, _ = lo.Difference(targetUrls, urls)
	if len(blackList) > 0 {
		_ = os.Truncate("blacklist.txt", 0)
	}

	body := fmt.Sprintf("%s\n", strings.Join(blackList, "\n"))
	_ = ioutil.WriteFile("blacklist.txt", []byte(body), os.ModePerm)

	log.Infof("target url len %d", len(targetUrls))
	log.Infof("url len %d", len(urls))
	log.Infof("black url len %d", len(blackList))
}

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

	urlsLen := len(urls)
	randomIndex := rand.Intn(urlsLen)
	if randomIndex >= urlsLen {
		return urls[0]
	} else {
		return urls[randomIndex]
	}
}

var client = &http.Client{
	Timeout: 3 * time.Second,
}

func fetchUrls(wg *sync.WaitGroup, k int, url string) {
	defer wg.Done()

	resp, err := client.Get(url)
	if err != nil {
		log.Errorf("fetch urls error: %s", err.Error())
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
