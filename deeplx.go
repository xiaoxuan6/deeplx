package deeplx

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/abadojack/whatlanggo"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"github.com/xiaoxuan6/deeplx/api/log"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
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

func fetchUri() string {
	if len(targetUrls) < 1 {
		client := &http.Client{
			Timeout: 3 * time.Second,
		}

		resp, err := client.Get("https://github-mirror.us.kg/https://github.com/ycvk/deeplx-local/blob/windows/url.txt")
		defer func() {
			_ = resp.Body.Close()
		}()

		if err == nil {
			r := bufio.NewReader(resp.Body)
			for {
				line, _, errs := r.ReadLine()
				if errs == io.EOF {
					break
				}

				targetUrls = append(targetUrls, string(line))
			}
			log.Infof("fetch urls len: %s", strconv.Itoa(len(targetUrls)))
		} else {
			log.Errorf("fetch urls error: %s", err.Error())
		}
	}

	urls = append(urls, targetUrls...)
	randomIndex := rand.Intn(len(urls))
	return urls[randomIndex]
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
	_ = retry.Do(
		func() error {
			uri := fetchUri()
			response, err := http.Post(uri, "application/json", strings.NewReader(string(jsonBody)))
			log.Info(fmt.Sprintf("url：%s, params：%s", uri, string(jsonBody)))

			if err == nil {
				defer func() {
					_ = response.Body.Close()
				}()

				body, err = io.ReadAll(response.Body)
				log.Infof("response：%s", string(body))
			} else {
				body = []byte(`{"code":500, "message": ` + err.Error() + `}`)
				log.Errorf("response error: %s", err.Error())
			}

			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	return &Response{
		Code: gjson.Get(string(body), "code").Int(),
		Data: gjson.Get(string(body), "data").String(),
		Msg:  gjson.Get(string(body), "message").String(),
	}
}
