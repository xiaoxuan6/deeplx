package deeplx

import (
	"bytes"
	"fmt"
	"github.com/OwO-Network/DeepLX/translate"
	"github.com/abadojack/whatlanggo"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"github.com/xiaoxuan6/deeplx/api/log"
	"io"
	"slices"
	"strings"
	"sync"
)

var (
	wg sync.WaitGroup

	blackList  = make([]string, 0)
	targetUrls = make([]string, 0)
	urls       = []string{"https://deeplx.mingming.dev/translate"}
)

type Response struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
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

	req := bytes.Buffer{}
	req.WriteString(`{"text":"` + text + `", "source_lang":"` + sourceLang + `", "target_lang":"` + targetLang + `"}`)

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

			response, err := client.Post(uri, "application/json", strings.NewReader(req.String()))
			log.Info(fmt.Sprintf("url：%s, params：%s", uri, req.String()))
			defer response.Body.Close()

			if err == nil {
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
