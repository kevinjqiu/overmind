compile:
	CGO_ENABLED=0 go build -tags netgo -a -v -o overmind cmd/main.go

build: compile
	docker build -t overmind .
