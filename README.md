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
	result := deeplx.Translate("Hello", "EN", "ZH")
	fmt.Println(result)
}
```

## 参数

|字段|描述|
|:--|:--|
|text|需要翻译的内容|
|source_lange|需要翻译语言|
|target_lange|目标语言|

## 自定义理由

默认路由 `translate`，指定变量 `ROUTER_PATH="translate"`

```docker
docker run --name=deeplx -e ROUTER_PATH="translate" -p 8311:8311 -d ghcr.io/xiaoxuan6/deeplx:latest 
```

## 允许打印日志，写入到文件中

```docker
docker run --name=deeplx -e VERBOSE="true" -p 8311:8311 -d ghcr.io/xiaoxuan6/deeplx:latest 
```

## 支持 `.env` 文件

文件内容为

```env
ROUTER_PATH="translate"
VERBOSE="true"
```