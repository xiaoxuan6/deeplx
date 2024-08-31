package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/xiaoxuan6/deeplx"
	"net/http"
)

type Response struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func SuccessResponse() string {
	res := &Response{
		Code: 200,
		Msg:  "welcome to deeplx translate",
	}

	b, _ := json.Marshal(res)
	return string(b)
}

func ErrorResponse(message string) string {
	res := &Response{
		Code: 500,
		Msg:  message,
	}

	b, _ := json.Marshal(res)
	return string(b)
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(ErrorResponse(fmt.Sprintf("method [%s] not allow", r.Method))))
}

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(ErrorResponse(fmt.Sprintf("route [%s] not found", r.URL.Path))))
}

type Request struct {
	Text       string `json:"text"`
	SourceLang string `json:"source_lang"`
	TargetLang string `json:"target_lang"`
}

func Translate(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var data Request
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	result := deeplx.Translate(data.Text, data.SourceLang, data.TargetLang)

	b, _ := json.Marshal(result)
	_, _ = w.Write(b)
}

func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(SuccessResponse()))
}
