package deeplx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate(t *testing.T) {
	response := Translate("Hello", "EN", "ZH")
	assert.Equal(t, int64(200), response.Code)
	t.Log(response)
}
