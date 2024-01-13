.PHONY: clean deps simplify run test build

clean:
		rm -rf target; \
		rm -f coverage.*

deps: clean
		go get -d -v ./...

simplify:
		gofmt -s -l -w .

test: deps
		go test -count=1 -v ./...

build: deps
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
		go build \
		-a -installsuffix cgo \
		-tags=jsoniter -o target/app .

run:
	go run main.go api

lint:
	docker run -t --rm -v $(PWD):/app -v ~/.cache/golangci-lint/v1.53.3:/root/.cache -w /app golangci/golangci-lint:v1.53.3 golangci-lint run
