SHELL := /bin/zsh
PROJECT_DIR=$(shell pwd)

gen-frontend-client:
	cd web && pnpm run gen

install-toolchains:
	$(PROJECT_DIR)/.anchor/bin/anchor install --config dev/anchor.yaml .

anchor-gen: install-toolchains
	$(PROJECT_DIR)/.anchor/bin/anchor gen --config dev/anchor.yaml .

###################################################
### Common
###################################################

gen: anchor-gen
	@go mod tidy

###################################################
### Dev enviornment
###################################################

start:
	docker compose up 

reload:
	docker compose restart dev

log:
	docker compose logs -f dev

db:
	psql "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

###################################################
### Build
###################################################

VERSION=v0.3.2

build-web:
	@cd web && pnpm install && pnpm run build

build-binary:
	@rm -rf upload
	@CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Darwin/x86_64/wavekit cmd/wavekit/main.go
	@CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Darwin/arm64/wavekit cmd/wavekit/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/x86_64/wavekit cmd/wavekit/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=386   go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/i386/wavekit cmd/wavekit/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/arm64/wavekit cmd/wavekit/main.go

binary-push:
	@cp scripts/download.sh upload/download.sh
	@echo 'latest version: $(VERSION)' > upload/metadata.txt
	@aws s3 cp --recursive upload/ s3://wavekit-release/	

build-server:
	GOOS=linux GOARCH=amd64 go build -o ./bin/wavekit-server-amd64 cmd/wavekit/main.go
	GOOS=linux GOARCH=arm64 go build -o ./bin/wavekit-server-arm64 cmd/wavekit/main.go

IMG_TAG=$(VERSION)
DOCKER_REPO=risingwavelabs/wavekit

push-docker: build-server
	docker buildx build --platform linux/amd64,linux/arm64 -f docker/Dockerfile.pgbundle -t ${DOCKER_REPO}:${IMG_TAG}-pgbundle --push .
	docker buildx build --platform linux/amd64,linux/arm64 -f docker/Dockerfile -t ${DOCKER_REPO}:${IMG_TAG} --push .

ci: doc build-web build-server build-binary push-docker binary-push

ut:
	@COLOR=ALWAYS go test -race -covermode=atomic -coverprofile=coverage.out -tags ut ./... 
	@grep -vE "_gen\.go|/mock[s]?/" coverage.out > coverage.filtered
	@go tool cover -func=coverage.filtered | fgrep total | awk '{print "Coverage:", $$3}'
	@go tool cover -html=coverage.filtered -o coverage.html


# https://pkg.go.dev/net/http/pprof#hdr-Usage_examples
pprof:
	go tool pprof http://localhost:8777/debug/pprof/$(ARG)
