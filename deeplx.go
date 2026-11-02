package deeplx

import (
	"bytes"
	"encoding/json"
	"github.com/avast/retry-go"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
	"strings"
)

type Response struct {
	Code int64  `json:"code"`
	Data string `json:"data"`
	Msg  string `json:"msg"`
}

// Translate
// ar 阿拉伯语
// de 德语
// ru 俄语
// fr 法语
// ko 韩语
// pt 葡萄牙语
// ja 日语
// th 泰语
// tr 土耳其语
// es 西班牙语
// it 意大利语
// en 英语
// vi 越南语
// ct 粤语
// zh 中文
func Translate(text string, languages ...string) *Response {
	sourceLang := "en"
	targetLang := "zh"

	if len(languages) == 1 {
		targetLang = strings.ToLower(languages[0])
	}

	if len(languages) == 2 && languages[0] != "" && languages[1] != "" {
		sourceLang = strings.ToLower(languages[0])
		targetLang = strings.ToLower(languages[1])
	}

	var body []byte
	err := retry.Do(
		func() error {
			cookies := []*http.Cookie{
				{Name: "pgv_info", Value: "ssid=s3205783488"},
				{Name: "pgv_pvid", Value: "941850630"},
				{Name: "vversion_name", Value: "8.2.95"},
				{Name: "_qpsvr_localtk", Value: "0.5439582279352217"},
				{Name: "video_omgid", Value: "2ce542b4a6241bc9"},
				{Name: "_ga", Value: "GA1.1.1025310258.1688290262"},
				{Name: "fqm_sessionid", Value: "06805962-6925-4a76-8ac8-efe3a848c0bc"},
				{Name: "qy_register_type", Value: "free"},
				{Name: "verifysession", Value: "h010bab6b15c66e4ec87a87c91a048ce066b139919a9add4d4904d747b19cd06d9905a8d6bb329b23ba"},
				{Name: "uin", Value: "o1527736751"},
				{Name: "_t_qbtool_uid", Value: "aaaanm6bx036fdsoz8vokozbwsqp88cb"},
				{Name: "_ga_TPFW0KPXC1", Value: "GS1.1.1720311535.24.1.1720311560.0.0.0"},
				{Name: "fqm_pvqid", Value: "0e0a0b58-20c4-467e-b299-9c2fe68389cf"},
				{Name: "eas_sid", Value: "m1s7h3w6i8s4q9g3Q7A1t7N3x8"},
				{Name: "101qqcomrouteLine", Value: "main_main_main"},
				{Name: "TSMT_CLIENT_KEY", Value: "browser-chrome-131.0.0-Windows_10-a2deca1f-355a-4b6a-850d-44578669a9a1-1737258172940"},
			}

			payload := map[string]interface{}{
				"header": map[string]string{
					"fn":         "auto_translation",
					"session":    "",
					"client_key": "browser-chrome-131.0.0-Windows_10-a2deca1f-355a-4b6a-850d-44578669a9a1-1737258172940",
					"user":       "",
				},
				"type":           "plain",
				"model_category": "normal",
				"text_domain":    "general",
				"source": map[string]interface{}{
					"lang":      sourceLang,
					"text_list": []string{text},
				},
				"target": map[string]string{
					"lang": targetLang,
				},
			}

			bodyBytes, err := json.Marshal(payload)
			if err != nil {
				panic(err)
			}

			req, err := http.NewRequest(http.MethodPost, "https://transmart.qq.com/api/imt", bytes.NewReader(bodyBytes))
			if err != nil {
				panic(err)
			}

			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/131.0.0.0 Safari/537.36")

			for _, c := range cookies {
				req.AddCookie(c)
			}

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				panic(err)
			}

			defer func() {
				_ = resp.Body.Close()
			}()

			body, err = io.ReadAll(resp.Body)
			return err
		},
		retry.Attempts(3),
		retry.LastErrorOnly(true),
	)

	if err == nil || gjson.GetBytes(body, "auto_translation.0").String() != "" {
		return &Response{
			Code: http.StatusOK,
			Data: gjson.Get(string(body), "auto_translation.0").String(),
			Msg:  "ok",
		}
	}

	return &Response{
		Code: http.StatusBadRequest,
		Data: "",
		Msg:  err.Error(),
	}
}
