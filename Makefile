build: login
	docker build --platform linux/amd64,linux/arm64 -t ghcr.io/xiaoxuan6/deeplx:latest . --push

login:
	@awk -F'[@:]' '/@github.com/ {print $$3}' /mnt/c/Users/Administrator/.git-credentials > token.txt
	@echo "$(shell cat token.txt)" | docker login ghcr.io --username xiaoxuan6 --password-stdin
	@rm -rf token.txt

