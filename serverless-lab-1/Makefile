.PHONY: build clean deploy

build:
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api api/main.go
	env GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api2 api2/main.go
	env GGOARC=amamd GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -o bin/api3 api3/main.go 

clean:
	rm -rf ./bin

deploy: clean build
	sls deploy --verbose
