run:
	go run github.com/rafaelsq/wtc

build-docker:
	CGO_ENABLED=0 go build -o api -ldflags="-s -w" cmd/api/api.go
	docker build -t rodrigobrito/rinha-api -f build/api.Dockerfile .
	CGO_ENABLED=0 go build -o storage -ldflags="-s -w" cmd/storage/storage.go
	docker build -t rodrigobrito/rinha-storage -f build/storage.Dockerfile .
