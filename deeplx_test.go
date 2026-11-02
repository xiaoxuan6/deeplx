package deeplx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate(t *testing.T) {
	response := Translate("Hello")
	assert.Equal(t, int64(200), response.Code)
	assert.Equal(t, "你好", response.Data)
	assert.Equal(t, "ok", response.Msg)
	t.Log(response)
}

func TestTranslateByLanguage(t *testing.T) {
	response := Translate("Hello", "en", "vi")
	assert.Equal(t, int64(200), response.Code)
	t.Log(response)

	response = Translate("Hello", "en", "ja")
	assert.Equal(t, int64(200), response.Code)
	t.Log(response)

	response = Translate("Hello", "th")
	assert.Equal(t, int64(200), response.Code)
	t.Log(response)
}
