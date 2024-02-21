package deeplx

import (
	"encoding/json"
	"errors"
	"github.com/abadojack/whatlanggo"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

func Translate(text, sourceLang, targetLang string) (string, error) {
	if len(sourceLang) == 0 {
		lang := whatlanggo.DetectLang(text)
		deepLLang := strings.ToUpper(lang.Iso6391())
		sourceLang = deepLLang
	}

	if len(targetLang) == 0 {
		targetLang = "EN"
	}

	if len(text) == 0 {
		return "", errors.New("No Translate Text Found")
	}

	response, err := post("https://api.deeplx.org/translate", RequestParams{
		Text:       text,
		SourceLang: sourceLang,
		TargetLang: targetLang,
	})

	if err != nil {
		return "", err
	}

	if gjson.Get(string(response), "code").Int() != 200 {
		return "", errors.New(gjson.Get(string(response), "message").String())
	}

	return gjson.Get(string(response), "data").String(), nil
}

type RequestParams struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

func post(url string, request RequestParams) ([]byte, error) {

	jsonBody, _ := json.Marshal(request)
	params := strings.NewReader(string(jsonBody))

	var body []byte
	err := retry.Do(
		func() error {
			response, err := http.Post(url, "application/json", params)

			if err == nil {
				defer func() {
					_ = response.Body.Close()
				}()

				body, err = io.ReadAll(response.Body)
			}

			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	return body, err
}
