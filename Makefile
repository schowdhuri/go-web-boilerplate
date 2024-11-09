.PHONY: build clean run

build:
	go run cmd/bundler.go
	go build -o bin/main main.go

clean:
	rm -rf bin/*
	rm -rf tmp/*
	rm -rf dist/*

run:
	./bin/main
