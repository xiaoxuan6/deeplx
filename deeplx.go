package deeplx

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/OwO-Network/DeepLX/translate"
	"github.com/abadojack/whatlanggo"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"github.com/xiaoxuan6/deeplx/api/log"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"slices"
	"strings"
	"sync"
	"time"
)

var (
	wg   sync.WaitGroup
	lock sync.Mutex

	blackList  = make([]string, 0)
	targetUrls = make([]string, 0)
	urls       = []string{"https://deeplx.mingming.dev/translate"}
)

type request struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

type Response struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

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
	_ = fetchUri()

	blackList = blackList[:0]
	for _, url := range targetUrls {
		url := url
		wg.Add(1)
		go func() {
			defer wg.Done()

			client := &http.Client{
				Timeout: 3 * time.Second,
			}
			resp, err := client.Get(strings.ReplaceAll(url, "/translate", ""))
			if err != nil {
				lock.Lock()
				blackList = append(blackList, url)
				lock.Unlock()
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				lock.Lock()
				blackList = append(blackList, url)
				lock.Unlock()
			}
		}()
	}

	wg.Wait()

	body := fmt.Sprintf("%s\n", strings.Join(blackList, "\n"))
	_ = ioutil.WriteFile("blacklist.txt", []byte(body), os.ModePerm)
}

func fetchUri() string {
	if len(targetUrls) < 1 {
		client := &http.Client{
			Timeout: 3 * time.Second,
		}

		resp, err := client.Get("https://github-mirror.us.kg/https://github.com/ycvk/deeplx-local/blob/windows/url.txt")
		if err != nil {
			log.Errorf("fetch urls error: %s", err.Error())
		} else {
			defer resp.Body.Close()

			r := bufio.NewReader(resp.Body)
			for {
				line, _, errs := r.ReadLine()
				if errs == io.EOF {
					break
				}

				targetUrls = append(targetUrls, string(line))
			}
		}

		for _, url := range targetUrls {
			wg.Add(1)
			url := url

			go func() {
				defer wg.Done()
				resp, err = client.Get(url)
				if err != nil {
					return
				}
				defer resp.Body.Close()

				if resp.StatusCode == 200 {
					urls = append(urls, url)
				}
			}()
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

func Translate(text, sourceLang, targetLang string) *Response {
	if len(text) == 0 {
		return &Response{
			Code: 500,
			Msg:  "No Translate Text Found",
		}
	}

	if len(sourceLang) == 0 {
		lang := whatlanggo.DetectLang(text)
		sourceLang = strings.ToUpper(lang.Iso6391())
	}

	if len(targetLang) == 0 {
		targetLang = "EN"
	}

	req := &request{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}
	jsonBody, _ := json.Marshal(req)

	var body []byte
	err := retry.Do(
		func() error {
			var uri string
			for {
				uri = fetchUri()
				if ok := slices.Contains(blackList, uri); !ok {
					break
				}
			}

			client := &http.Client{
				Timeout: 3 * time.Second,
			}
			response, err := client.Post(uri, "application/json", strings.NewReader(string(jsonBody)))
			log.Info(fmt.Sprintf("url：%s, params：%s", uri, string(jsonBody)))

			if err == nil {
				defer func() {
					_ = response.Body.Close()
				}()

				body, err = io.ReadAll(response.Body)
				log.Infof("response：%s", string(body))
			} else {
				blackList = append(blackList, uri)
				body = []byte(`{"code":500, "message": ` + err.Error() + `}`)
				log.Errorf("response error: %s", err.Error())
			}

			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	if err == nil {
		return &Response{
			Code: gjson.Get(string(body), "code").Int(),
			Data: gjson.Get(string(body), "data").String(),
			Msg:  gjson.Get(string(body), "message").String(),
		}
	}

	result, err := translate.TranslateByDeepLX(sourceLang, targetLang, text, "", "")
	if err != nil {
		return &Response{
			Code: 500,
			Msg:  err.Error(),
		}
	}

	return &Response{
		Code: int64(result.Code),
		Data: result.Data,
		Msg:  result.Message,
	}
}
