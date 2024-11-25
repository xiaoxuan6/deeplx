package deeplx

import (
	"bytes"
	"github.com/OwO-Network/DeepLX/translate"
	"github.com/abadojack/whatlanggo"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup

	targetUrls = make([]string, 0)
	urls       = []string{"https://deeplx.mingming.dev/translate"}
)

type Response struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

func Translate(text, sourceLang, targetLang string) *Response {
	return TranslateWithProxyUrl(text, sourceLang, targetLang, false)
}

func TranslateWithProxyUrl(text, sourceLang, targetLang string, isProxyUrl bool) *Response {
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

	requestParams := bytes.Buffer{}
	requestParams.WriteString(`{"text":"` + text + `","source_lang":"` + sourceLang + `","target_lang":"` + targetLang + `"}`)

	transport := &http.Transport{}
	if isProxyUrl == true {
		if proxyUrl := getProxyUrl(); proxyUrl != "" {
			proxy, errs := url.Parse(proxyUrl)
			if errs == nil {
				transport.Proxy = http.ProxyURL(proxy)
			}
		}
	}

	httpClient := &http.Client{
		Transport: transport,
		Timeout:   3 * time.Second,
	}

	var body []byte
	err := retry.Do(
		func() error {
			response, err := httpClient.Post(fetchUri(), "application/json", strings.NewReader(requestParams.String()))
			if err != nil {
				body = []byte(`{"code":500, "message": "` + err.Error() + `"}`)
				return err
			}

			defer func() {
				_ = response.Body.Close()
			}()

			body, err = io.ReadAll(response.Body)
			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	if err != nil {
		var proxyUrl string
		if isProxyUrl == true {
			proxyUrl = getProxyUrl()
		}
		return TranslateByDeeplx(text, sourceLang, targetLang, proxyUrl)
	}

	return &Response{
		Code: gjson.Get(string(body), "code").Int(),
		Data: gjson.Get(string(body), "data").String(),
		Msg:  gjson.Get(string(body), "message").String(),
	}
}

func TranslateByDeeplx(text, sourceLang, targetLang, proxyUrl string) *Response {
	result, err := translate.TranslateByDeepLX(sourceLang, targetLang, text, "", proxyUrl)
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
