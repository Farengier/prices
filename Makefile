build:
	go build -o build/prices src/cli/prices.go

run:
	build/prices -c example/config.yml

.PHONY: build run