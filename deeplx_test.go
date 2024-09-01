package deeplx

import (
	"github.com/stretchr/testify/assert"
	"github.com/xiaoxuan6/deeplx/api/log"
	"os"
	"testing"
)

func TestTranslate(t *testing.T) {
	_ = os.Setenv("VERBOSE", "true")
	log.InitLog()

	response := Translate("Hello", "EN", "ZH")
	assert.Equal(t, int64(200), response.Code)
	t.Log(response)
}

func TestCheckUrlAndReloadBlack(t *testing.T) {
	CheckUrlAndReloadBlack()
	assert.Equal(t, nil, "")
}
