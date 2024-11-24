build:
	docker build -t ghcr.io/xiaoxuan6/deeplx:latest .

docker-login:
	@awk -F'[@:]' '/@github.com/ {print $$3}' /mnt/c/Users/Administrator/.git-credentials > token.txt
	@echo "$(shell cat token.txt)" | docker login ghcr.io --username xiaoxuan6 --password-stdin
	@rm -rf token.txt

push:
	docker push ghcr.io/xiaoxuan6/deeplx:latest

