TAG=nobonobo/webpush-demo

.PHONY: all build generate run pub

all: generate build

build:
	go build .

generate:
	go generate .

run: build
	./webpush-demo

pub:
	cloudflared tunnel --url http://localhost:8080

docker:
	docker build --rm -t $(TAG) .
	docker login
	docker push $(TAG)
