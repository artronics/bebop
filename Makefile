.SILENT:

-include .env

token:
	get_token -p $(username)

build:
	go build -o build/bebop

run: build
	./build/bebop $(opt)

dev: build
	./build/bebop $(opt_dev)

.PHONY: build run
