.PHONY: clean
clean:
	rm overmind

overmind:
	CGO_ENABLED=0 go build -tags netgo -a -v -o overmind cmd/main.go

.PHONY: image
image: overmind
	docker build -t overmind .

.PHONY: bootstrap-image
bootstrap-image:
	docker build -f Dockerfile.bootstrap -t overmind-bootstrap .