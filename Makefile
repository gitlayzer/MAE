.PHONY:

tag: $(shell git rev-parse --short HEAD)
hub_user: layzer
release_name: mae

run: build
	./mae run

build:
	go build -o mae main.go

image:
	docker build -t mae:$(tag) .

push:
	docker push $(hub_user)/$(release_name):$(tag)



