# deeplx
A Go library used for unlimited DeepL translation

## Installation

Install it with the go get command:
```bash
go get github.com/xiaoxuan6/deeplx
```

## Usage
Then, you can create a new DeepL translation client and use it for translation:

```go
import (
	"fmt"
	"github.com/xiaoxuan6/deeplx"
)

func main() {
	result, err := deeplx.Translate("Hello", "EN", "ZH")
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Println(result)
}
```