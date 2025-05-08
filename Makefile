SHELL := /bin/zsh
PROJECT_DIR=$(shell pwd)

gen-frontend-client:
	cd web && pnpm run gen

install-anchor:
	@GOBIN=$(PROJECT_DIR)/bin go install github.com/cloudcarver/anchor/cmd/anchor@latest

anchor-gen: install-anchor
	$(PROJECT_DIR)/bin/anchor gen .

###################################################
### Common
###################################################

gen: anchor-gen doc gen-frontend-client
	@go mod tidy

###################################################
### Documentation
###################################################

CONFTEXT_VERSION=v0.3.3
CONFTEXT_BIN=$(PROJECT_DIR)/bin/conftext

install-doc-tools:
	@GOBIN=$(PROJECT_DIR)/bin BIN=conftext VERSION=${CONFTEXT_VERSION} DIR=$(PROJECT_DIR)/bin REPO=github.com/cloudcarver/edc/cmd/conftext ./scripts/go-install.sh

doc-config:
	@awk -v cmds='$(CONFTEXT_BIN) -prefix wk -path internal/config -yaml|CONFIG_SAMPLE_YAML;\
		$(CONFTEXT_BIN) -prefix wk -path internal/config -env -markdown|CONFIG_ENV;\
		cat dev/init.yaml|CONFIG_SAMPLE_INIT' \
		-f scripts/template-subst.awk docs/templates/config.tmpl.md > docs/config.md

doc-contributing:
	@awk -v cmds='cat CONTRIBUTING.md|CONTRIBUTING_MD' \
		-f scripts/template-subst.awk docs/templates/CONTRIBUTING.tmpl.md > CONTRIBUTING.md

doc: install-doc-tools doc-config doc-contributing

###################################################
### Dev enviornment
###################################################

K0S_KUBECTL=docker exec -ti wavekit-k0s k0s kubectl
K0S_CODEBASE_DIR=/opt/wavekit-dev/codebase

start:
	docker-compose up -d
	./dev/init.sh
	$(K0S_KUBECTL) apply -f $(K0S_CODEBASE_DIR)/dev/k0s.yaml > /dev/null 2>&1

apply:
	$(K0S_KUBECTL) apply -f $(K0S_CODEBASE_DIR)/dev/k0s.yaml

reload:
	$(K0S_KUBECTL) rollout restart deployment/wavekit

log:
	$(K0S_KUBECTL) logs -l app=wavekit --follow

db:
	psql "postgresql://postgres:postgres@localhost:30432/postgres?sslmode=disable"

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
