package deeplx

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTranslate(t *testing.T) {
	response, err := Translate("Hello", "EN", "ZH")
	assert.Nil(t, err)
	t.Log(response)
}
