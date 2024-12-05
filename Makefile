all: build delete_data
deploy: build gofumpt delete_data

build:
	go build -o s3 .

gofumpt:
	gofumpt -l -w .

delete_data:
	rm -rf data