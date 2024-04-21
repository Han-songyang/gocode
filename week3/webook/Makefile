.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=amd64 go build -o webook .
	@docker rmi -f gocode/webook:v0.1
	@docker build -t han0426/webook:v0.2 .
