all:
	- make install
	- make test

install:
	- go get -v ./...
	- bower install
	- mkdir keys
	- openssl genrsa -out keys/app.rsa 1024
	- openssl rsa -in keys/app.rsa -pubout > keys/app.rsa.pub

test:
	- go test -v ./models ./jarvis

build:
	- rm steel
	- go build

run:
	- make build
	- ./steel
