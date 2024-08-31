# deeplx
A Go library used for unlimited DeepL translation

## 参数

|字段|描述|
|:--|:--|
|text|需要翻译的内容|
|source_lang|需要翻译语言|
|target_lang|目标语言|

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

## 支持黑名单

在文件 `blacklist.txt` 中添加不可用的 `url` 地址

```docker
docker run --name=deeplx -v $(pwd)/blacklist.txt:/blacklist.txt -p 8311:8311 -d ghcr.io/xiaoxuan6/deeplx:latest 
```