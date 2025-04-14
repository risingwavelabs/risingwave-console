SHELL := /bin/zsh
PROJECT_DIR=$(shell pwd)

EE ?= false

###################################################
### OpenAPI         
###################################################

OAPI_CODEGEN_VERSION=v2.4.1
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
ifeq ($(EE), true)
	$(WIRE_GEN) ./ee/wire
else
	$(WIRE_GEN) ./wire
endif

###################################################
### SQL  
###################################################

SQLC_VERSION=v1.27.0
QUERIER_DIR=$(PROJECT_DIR)/internal/model/querier
SQLC_BIN=$(PROJECT_DIR)/bin/sqlc

install-sqlc:
	@DIR=$(PROJECT_DIR)/bin VERSION=${SQLC_VERSION} ./scripts/install-sqlc.sh

clean-querier:
	@rm -f $(QUERIER_DIR)/*.sql.gen.go || true
	@rm -f $(QUERIER_DIR)/copyfrom_gen.go   
	@rm -f $(QUERIER_DIR)/db_gen.go
	@rm -f $(QUERIER_DIR)/models_gen.go
	@rm -f $(QUERIER_DIR)/querier_gen.go

gen-querier: install-sqlc clean-querier
	$(SQLC_BIN) generate

###################################################
### mock 
###################################################

MOCKGEN_VERSION=0.5.0
MOCKGEN_BIN=$(PROJECT_DIR)/bin/mockgen

install-mockgen: 
	@DIR=$(PROJECT_DIR)/bin VERSION=${MOCKGEN_VERSION} ./scripts/install-mockgen.sh

gen-mock: install-mockgen
	$(MOCKGEN_BIN) -source=internal/model/model.go -destination=internal/model/mock_gen.go -package=model
	$(MOCKGEN_BIN) -source=internal/task/task.go -destination=internal/task/mock/task_mock_gen.go -package=mock
	$(MOCKGEN_BIN) -source=internal/worker/lifecycle_handler.go -destination=internal/worker/mock/lifecycle_handler_mock_gen.go -package=mock
	$(MOCKGEN_BIN) -source=internal/worker/worker.go -destination=internal/worker/mock/worker_mock_gen.go -package=mock
	$(MOCKGEN_BIN) -source=internal/service/service.go -destination=internal/service/service_mock_gen.go -package=service
	$(MOCKGEN_BIN) -source=internal/modelctx/modelctx.go -destination=internal/modelctx/mock/modelctx_mock_gen.go -package=mock
	$(MOCKGEN_BIN) -source=internal/conn/meta/types.go -destination=internal/conn/meta/mock/mock_gen.go -package=mock
	$(MOCKGEN_BIN) -source=internal/task/task.go -destination=internal/task/task_mock_gen.go -package=task
	$(MOCKGEN_BIN) -source=internal/conn/http/http.go -destination=internal/conn/http/mock/http_mock_gen.go -package=mock

###################################################
### Common
###################################################

gen: doc gen-spec gen-querier gen-wire gen-mock gen-frontend-client
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
### Promdump
###################################################

upload-promdump:
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o upload/promdump/Darwin/x86_64/promdump cmd/promdump/main.go
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o upload/promdump/Darwin/arm64/promdump cmd/promdump/main.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o upload/promdump/Linux/x86_64/promdump cmd/promdump/main.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=386   go build -o upload/promdump/Linux/i386/promdump cmd/promdump/main.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o upload/promdump/Linux/arm64/promdump cmd/promdump/main.go
	chmod +x upload/promdump/Darwin/x86_64/promdump
	chmod +x upload/promdump/Darwin/arm64/promdump
	chmod +x upload/promdump/Linux/x86_64/promdump
	chmod +x upload/promdump/Linux/i386/promdump
	chmod +x upload/promdump/Linux/arm64/promdump
	cp scripts/download-promdump.sh upload/promdump/download.sh
	aws s3 cp --recursive upload/promdump/ s3://wavekit-release/promdump/

###################################################
### Prompush
###################################################

upload-prompush:
	CGO_ENABLED=0 GOOS=darwin  GOARCH=amd64 go build -o upload/prompush/Darwin/x86_64/prompush cmd/prompush/main.go
	CGO_ENABLED=0 GOOS=darwin  GOARCH=arm64 go build -o upload/prompush/Darwin/arm64/prompush cmd/prompush/main.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=amd64 go build -o upload/prompush/Linux/x86_64/prompush cmd/prompush/main.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=386   go build -o upload/prompush/Linux/i386/prompush cmd/prompush/main.go
	CGO_ENABLED=0 GOOS=linux   GOARCH=arm64 go build -o upload/prompush/Linux/arm64/prompush cmd/prompush/main.go
	chmod +x upload/prompush/Darwin/x86_64/prompush
	chmod +x upload/prompush/Darwin/arm64/prompush
	chmod +x upload/prompush/Linux/x86_64/prompush
	chmod +x upload/prompush/Linux/i386/prompush
	chmod +x upload/prompush/Linux/arm64/prompush
	cp scripts/download-prompush.sh upload/prompush/download.sh
	aws s3 cp --recursive upload/prompush/ s3://wavekit-release/prompush/

###################################################
### Dev enviornment
###################################################

K0S_KUBECTL=docker exec -ti wavekit-k0s k0s kubectl
K0S_CODEBASE_DIR=/opt/wavekit-dev/codebase

start:
	docker-compose up -d
	./dev/init.sh
	$(K0S_KUBECTL) apply -f $(K0S_CODEBASE_DIR)/dev/k0s.yaml

apply:
	$(K0S_KUBECTL) apply -f $(K0S_CODEBASE_DIR)/dev/k0s.yaml

reload:
	$(K0S_KUBECTL) rollout restart deployment/wavekit

log:
	$(K0S_KUBECTL) logs -l app=wavekit --follow

db:
	psql "postgresql://postgres:postgres@localhost:5432/postgres?sslmode=disable"

test:
	TEST_DIR=$(PROJECT_DIR)/e2e HOLD="$(HOLD)" ./scripts/run-local-test.sh "$(K)" 


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
	@go tool cover -html coverage.out -o coverage.html
	@go tool cover -func coverage.out | fgrep total | awk '{print "Coverage:", $$3}'


# https://pkg.go.dev/net/http/pprof#hdr-Usage_examples
pprof:
	go tool pprof http://localhost:8777/debug/pprof/$(ARG)
