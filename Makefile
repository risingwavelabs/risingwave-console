SHELL := /bin/zsh
PROJECT_DIR=$(shell pwd)

###################################################
### OpenAPI         
###################################################

OAPI_CODEGEN_VERSION=v2.1.0
OAPI_CODEGEN_BIN=$(PROJECT_DIR)/bin/oapi-codegen
OAPI_GEN_DIR=$(PROJECT_DIR)/internal/apigen
OAPI_CODEGEN_FIBER_BIN=$(PROJECT_DIR)/bin/oapi-codegen-fiber

install-oapi-codegen:
	@DIR=$(PROJECT_DIR)/bin VERSION=${OAPI_CODEGEN_VERSION} ./scripts/install-oapi-codegen.sh
	
install-oapi-codegen-fiber:
	@GOBIN=$(PROJECT_DIR)/bin go install github.com/cloudcarver/oapi-codegen-fiber@v0.5.1

prune-spec:
	@rm -f $(OAPI_GEN_DIR)/spec_gen.go

OAPI_GENERATE_ARG=types,fiber,client

gen-spec: install-oapi-codegen-fiber install-oapi-codegen prune-spec
	$(OAPI_CODEGEN_BIN) -generate $(OAPI_GENERATE_ARG) -o $(OAPI_GEN_DIR)/spec_gen.go -package apigen $(PROJECT_DIR)/web/api/v1.yaml
	$(PROJECT_DIR)/bin/oapi-codegen-fiber --package apigen --path $(PROJECT_DIR)/web/api/v1.yaml --out $(PROJECT_DIR)/internal/apigen/scopes_extend_gen.go

gen-frontend-client:
	cd web && pnpm run gen

###################################################
### Wire
###################################################

WIRE_VERSION=v0.6.0

install-wire:
	@DIR=$(PROJECT_DIR)/bin VERSION=${WIRE_VERSION} ./scripts/install-wire.sh

WIRE_GEN=$(PROJECT_DIR)/bin/wire
gen-wire: install-wire
	$(WIRE_GEN) ./wire

###################################################
### SQL  
###################################################

SQLC_VERSION=v1.27.0
QUERIER_DIR=$(PROJECT_DIR)/internal/model/querier
SQLC_BIN=$(PROJECT_DIR)/bin/sqlc

install-sqlc:
	@DIR=$(PROJECT_DIR)/bin VERSION=${SQLC_VERSION} ./scripts/install-sqlc.sh

clean-querier:
	@rm -f $(QUERIER_DIR)/*sql.gen.go
	@rm -f $(QUERIER_DIR)/copyfrom_gen.go   
	@rm -f $(QUERIER_DIR)/db_gen.go
	@rm -f $(QUERIER_DIR)/models_gen.go
	@rm -f $(QUERIER_DIR)/querier_gen.go

gen-querier: install-sqlc clean-querier
	$(SQLC_BIN) generate

###################################################
### mock 
###################################################

MOCKGEN_VERSION=1.6.0
MOCKGEN_BIN=$(PROJECT_DIR)/bin/mockgen

install-mockgen: 
	@DIR=$(PROJECT_DIR)/bin VERSION=${MOCKGEN_VERSION} ./scripts/install-mockgen.sh

gen-mock: install-mockgen
	$(MOCKGEN_BIN) -source=internal/model/model.go -destination=internal/model/mock_gen.go -package=model

###################################################
### Common
###################################################

gen: gen-spec gen-querier gen-wire gen-mock gen-frontend-client
	@go mod tidy

###################################################
### Documentation
###################################################

CONFTEXT_VERSION=v0.3.3
CONFTEXT_BIN=$(PROJECT_DIR)/bin/conftext

install-doc-tools:
	@GOBIN=$(PROJECT_DIR)/bin BIN=conftext VERSION=${CONFTEXT_VERSION} DIR=$(PROJECT_DIR)/bin REPO=github.com/cloudcarver/edc/cmd/conftext ./scripts/go-install.sh

doc-readme:
	@awk -v cmds='cat examples/docker-compose/docker-compose.yaml|README_DOCKER_COMPOSE' \
		-f scripts/template-subst.awk docs/templates/README.tmpl.md > README.md

doc-config:
	@awk -v cmds='$(CONFTEXT_BIN) -prefix wk -path internal/config -yaml|CONFIG_SAMPLE_YAML;\
		$(CONFTEXT_BIN) -prefix wk -path internal/config -env -markdown|CONFIG_ENV;\
		cat init.yaml|CONFIG_SAMPLE_INIT' \
		-f scripts/template-subst.awk docs/templates/config.tmpl.md > docs/config.md

doc-contributing:
	@awk -v cmds='cat CONTRIBUTING.md|CONTRIBUTING_MD' \
		-f scripts/template-subst.awk docs/templates/CONTRIBUTING.tmpl.md > CONTRIBUTING.md

doc: install-doc-tools doc-readme doc-config doc-contributing

###################################################
### Dev enviornment
###################################################

dev:
	docker-compose up

reload:
	docker-compose restart dev

db:
	psql "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

test:
	TEST_DIR=$(PROJECT_DIR)/e2e HOLD="$(HOLD)" ./scripts/run-local-test.sh "$(K)" 


###################################################
### Build
###################################################

VERSION=v0.1.2


build-web:
	@cd web && pnpm run build

release: build-web
	@rm -rf upload
	@CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Darwin/x86_64/wavekit cmd/main.go
	@CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Darwin/arm64/wavekit cmd/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/x86_64/wavekit cmd/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=386   go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/i386/wavekit cmd/main.go
	@CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -ldflags="-X 'github.com/risingwavelabs/wavekit/internal/utils.CurrentVersion=$(VERSION)'" -o upload/Linux/arm64/wavekit cmd/main.go
	@cp scripts/download.sh upload/download.sh
	@echo 'latest version: $(VERSION)' > upload/metadata.txt
	@aws s3 cp --recursive upload/ s3://wavekit-release/	

build-server:
	GOOS=linux GOARCH=amd64 go build -o ./bin/wavekit-server cmd/main.go

IMG_TAG=$(VERSION)
DOCKER_REPO=risingwavelabs/wavekit

build-docker:
	docker build -f docker/Dockerfile.pgbundle -t ${DOCKER_REPO}:${IMG_TAG}-pgbundle .
	docker build -f docker/Dockerfile -t ${DOCKER_REPO}:${IMG_TAG} .

docker-push:
	docker push ${DOCKER_REPO}:${IMG_TAG}-pgbundle
	docker push ${DOCKER_REPO}:${IMG_TAG}

build: build-web build-server build-docker

push: docker-push

ut:
	@COLOR=ALWAYS go test -race -covermode=atomic -coverprofile=coverage.out -tags ut ./... 
	@go tool cover -html coverage.out -o coverage.html
	@go tool cover -func coverage.out | fgrep total | awk '{print "Coverage:", $$3}'


