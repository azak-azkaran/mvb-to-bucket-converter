fetch:
	go mod download

build: fetch
	go build -o main

coverage: fetch
	@echo Running Test with Coverage export
	go test -v -coverprofile=cover.out
	go test -json > report.json
	go tool cover -html=cover.out -o cover.html


run: build
	./main $1
	cat ($1_converted.csv)
