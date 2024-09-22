.PHONY: build
build:
	go build .

.PHONY: debugger
debugger:
	dlv debug --headless --api-version=2 --listen=127.0.0.1:43000 .
