test:
	go test -v
build:
	# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags="-s -w" -o app -work $(shell pwd)/dev/
	docker build -t tttlkkkl/deploy $(shell pwd)/
	docker push tttlkkkl/deploy