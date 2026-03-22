
BINARY_NAME = mydocker
MAIN_PATH   = ./Docker/main.go

build:
	go build -o $(BINARY_NAME) $(MAIN_PATH)

run: build
	./$(BINARY_NAME) run -ti --name test busybox sh

clean:
	rm -rf $(BINARY_NAME)

tidy:
	go mod tidy

.PHONY: build run clean