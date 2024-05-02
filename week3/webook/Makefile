.PHONY: docker
docker:
	@rm webook || true
	@go mod tidy
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f han0426/webook:v0.2
	@docker build -t han0426/webook:v0.2 .
