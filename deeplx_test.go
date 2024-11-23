package deeplx

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestTranslate(t *testing.T) {
	response := Translate("Hello", "EN", "ZH")
	assert.Equal(t, int64(200), response.Code)
	t.Log(response)
}

func TestTranslateWithGo(t *testing.T) {
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			response := Translate("Hello", "EN", "ZH")
			assert.Equal(t, int64(200), response.Code)
			t.Log(response)
		}()
	}
	wg.Wait()

	end := time.Now().Sub(start).Seconds()
	t.Log(fmt.Sprintf("time: %.2f", end))
}

func TestTranslateByDeeplx(t *testing.T) {
	response := TranslateByDeeplx("Hello", "En", "zh")
	assert.Equal(t, int64(200), response.Code)
	t.Log(response)
}
