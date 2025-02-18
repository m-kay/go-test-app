
build: generate
	go build ./src/...
build-app: generate
	go build -o app ./src/main.go
run: build-app
	./app
test: build
	go test ./src/...
generate:
	go generate ./src/...