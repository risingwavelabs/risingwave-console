SHELL := /bin/zsh
PROJECT_DIR=$(shell pwd)
ANCHOR_BIN="$(PROJECT_DIR)/.anchor/bin/anchor"

gen-frontend-client:
	cd web && pnpm run gen

install-toolchains:
	$(ANCHOR_BIN) install --config dev/anchor.yaml .

anchor-gen: install-toolchains
	$(ANCHOR_BIN) gen --config dev/anchor.yaml .

###################################################
### Common
###################################################

gen: anchor-gen
	@go mod tidy

###################################################
### Documentation
###################################################

doc-config:
	@mkdir -p .tmp
	@$(ANCHOR_BIN) docs config --prefix RCONSOLE --path pkg/config --yaml > .tmp/sample_config.yaml
	@$(ANCHOR_BIN) docs config --prefix RCONSOLE --path pkg/config --env --markdown > .tmp/sample_config.env
	@cat docs/templates/config.tmpl.md | $(ANCHOR_BIN) docs replace --key CONFIG_SAMPLE_YAML --file .tmp/sample_config.yaml |\
		$(ANCHOR_BIN) docs replace --key CONFIG_ENV --file .tmp/sample_config.env |\
		$(ANCHOR_BIN) docs replace --key CONFIG_SAMPLE_INIT --file dev/init.yaml \
		> docs/config.md

doc: doc-config

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

VERSION=v0.4.0

build-web:
	@cd web && pnpm install && pnpm run build

build-binary:
	@rm -rf upload
	@CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags="-X 'github.com/risingwavelabs/risingwave-console/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Darwin/x86_64/risingwave-console cmd/risingwave-console/main.go
	@CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -ldflags="-X 'github.com/risingwavelabs/risingwave-console/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Darwin/arm64/risingwave-console cmd/risingwave-console/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags="-X 'github.com/risingwavelabs/risingwave-console/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/x86_64/risingwave-console cmd/risingwave-console/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=386   go build -ldflags="-X 'github.com/risingwavelabs/risingwave-console/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/i386/risingwave-console cmd/risingwave-console/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -ldflags="-X 'github.com/risingwavelabs/risingwave-console/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/arm64/risingwave-console cmd/risingwave-console/main.go

push-binary:
	@cp scripts/download.sh upload/download.sh
	@echo 'latest version: $(VERSION)' > upload/metadata.txt
	@aws s3 cp --recursive upload/ s3://risingwave-console/	

build-server:
	GOOS=linux GOARCH=amd64 go build -o ./bin/risingwave-console-server-amd64 cmd/risingwave-console/main.go
	GOOS=linux GOARCH=arm64 go build -o ./bin/risingwave-console-server-arm64 cmd/risingwave-console/main.go

IMG_TAG=$(VERSION)
DOCKER_REPO=risingwavelabs/risingwave-console

push-docker: build-server
	docker buildx build --platform linux/amd64,linux/arm64 -f docker/Dockerfile.pgbundle -t ${DOCKER_REPO}:${IMG_TAG}-pgbundle --push .
	docker buildx build --platform linux/amd64,linux/arm64 -f docker/Dockerfile -t ${DOCKER_REPO}:${IMG_TAG} --push .

ci: doc build-web build-server build-binary push-docker push-binary

ut:
	@COLOR=ALWAYS go test -race -covermode=atomic -coverprofile=coverage.out -tags ut ./... 
	@grep -vE "_gen\.go|/mock[s]?/" coverage.out > coverage.filtered
	@go tool cover -func=coverage.filtered | fgrep total | awk '{print "Coverage:", $$3}'
	@go tool cover -html=coverage.filtered -o coverage.html

prepare-test:
	cd test && uv sync
	cd test && uv run openapi-python-client generate --path ../api/v1.yaml --output-path oapi --overwrite 

test: prepare-test
	cd test && uv run main.py

# https://pkg.go.dev/net/http/pprof#hdr-Usage_examples
pprof:
	go tool pprof http://localhost:8777/debug/pprof/$(ARG)
