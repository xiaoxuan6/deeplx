package deeplx

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"net/http"
)

func main() {
	headers := map[string]string{
		"accept":             "application/json, text/plain, */*",
		"accept-language":    "zh-CN,zh;q=0.9,en;q=0.8",
		"cache-control":      "no-cache",
		"content-type":       "application/json",
		"origin":             "https://transmart.qq.com",
		"pragma":             "no-cache",
		"priority":           "u=1, i",
		"referer":            "https://transmart.qq.com/zh-CN/index",
		"sec-ch-ua":          "\"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"156\", \"Google Chrome\";v=\"156\"",
		"sec-ch-ua-mobile":   "?0",
		"sec-ch-ua-platform": "\"Windows\"",
		"sec-fetch-dest":     "empty",
		"sec-fetch-mode":     "cors",
		"sec-fetch-site":     "same-origin",
		"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/156.0.0.0 Safari/537.36",
		"x-requested-with":   "XMLHttpRequest",
	}

	cookies := map[string]string{
		"pgv_pvid":                    "941850630",
		"RK":                          "y8tEXYJfME",
		"ptcz":                        "394cd54c30442fc99a3ed996d0ba9b5c39c87d397fd9414cfa6a05fb4bd195ef",
		"qq_domain_video_guid_verify": "84ab424cc25a9f0a",
		"o_cookie":                    "1527736751",
		"_qimei_q36":                  "",
		"_qimei_h38":                  "5cd2bced712608f5400ef60802000008d19a10",
		"_qimei_i_1":                  "21c86f87c60e528992c2aa375a8675e0ffb8a1f5125b56d4b28f20582493206c61633fc13980b3ddd6b2c2d2",
		"a_sk__07__0WEB071JJOW4VSP7":  "015465636d64616c63606d6c",
		"TSMT_CLIENT_KEY":             "browser-chrome-156.0.0-Windows_10-ec6b0c90-f00d-4cef-a810-10975cd12f60-1793616655646",
	}

	url := "https://transmart.qq.com/api/imt"

	payload := map[string]interface{}{
		"header": map[string]string{
			"fn":         "support_lang",
			"session":    "",
			"client_key": "browser-chrome-156.0.0-Windows_10-ec6b0c90-f00d-4cef-a810-10975cd12f60-1793616655646",
			"user":       "",
		},
		"type":    "auto_translation",
		"sort_by": "pinyin",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		panic(err)
	}

	// 设置 headers
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	// 设置 cookies
	for k, v := range cookies {
		req.AddCookie(&http.Cookie{
			Name:  k,
			Value: v,
		})
	}

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(respBody))
	fmt.Println(resp.Status)

	codes := gjson.Get(string(respBody), "full_lang_pair.#.source.code").Array()
	chn_names := gjson.Get(string(respBody), "full_lang_pair.#.source.chn_name").Array()

	for i := 0; i < len(codes) && i < len(chn_names); i++ {
		println(codes[i].String(), chn_names[i].String())
	}
}
