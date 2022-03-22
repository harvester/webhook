.PHONY: docker, build, apply, clean

docker: build
	@docker build -t armarny/webhook-example:latest .
build:
	@GOOS=linux go build -o bin/webhook example/main.go
apply:
	@kubectl apply -f example/manifests
clean:
	@rm -rf bin
