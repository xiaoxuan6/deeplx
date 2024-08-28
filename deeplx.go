package deeplx

import (
	"bufio"
	"encoding/json"
	"github.com/abadojack/whatlanggo"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"io"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var urls = []string{"https://deeplx.mingming.dev/translate", "https://deeplx.niubipro.com/translate"}

type Request struct {
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

			urls = append(urls, string(line))
		}
	}

	randomIndex := rand.Intn(len(urls))
	return urls[randomIndex]
}

func Translate(text, sourceLang, targetLang string) Response {
	if len(text) == 0 {
		return Response{
			Code: 500,
			Msg:  "No Translate Text Found",
		}
	}

	if len(sourceLang) == 0 {
		lang := whatlanggo.DetectLang(text)
		deepLLang := strings.ToUpper(lang.Iso6391())
		sourceLang = deepLLang
	}

	if len(targetLang) == 0 {
		targetLang = "EN"
	}

	request := &Request{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	}
	jsonBody, _ := json.Marshal(request)

	var body []byte
	_ = retry.Do(
		func() error {
			response, err := http.Post(fetchUri(), "application/json", strings.NewReader(string(jsonBody)))

			if err == nil {
				defer func() {
					_ = response.Body.Close()
				}()

				body, err = io.ReadAll(response.Body)
			} else {
				body = []byte(`{"code":500, "message": ` + err.Error() + `}`)
			}

			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	return Response{
		Code: gjson.Get(string(body), "code").Int(),
		Data: gjson.Get(string(body), "data").String(),
		Msg:  gjson.Get(string(body), "message").String(),
	}
}
