.PHONY: docker
docker:
	@rm webook-app || true
	@GOOS=linux GOARCH=amd64 go build -tags=k8s -o webook-app .
	@docker rmi -f xjiang91/webook:v0.0.1
	@docker build -t xjiang91/webook:v0.0.1 .
