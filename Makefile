.PHONY: clean

all: build-all
# 子路径
SUB_DIR=./web

build-web:
	@make -C ${SUB_DIR} build
	@go generate main.go

build-all: build-web
	@mkdir -p ./_output
	@go build -o ./_output ./

dev: build-web
	@go run main.go

build-end:
	@mkdir -p ./_output
	@go build -o ./_output ./

clean:
	@rm -rf ./_output