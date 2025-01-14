PROJECT_DIR=$(shell pwd)

###################################################
### OpenAPI         
###################################################

OAPI_CODEGEN_VERSION=v2.1.0
OAPI_CODEGEN_BIN=$(PROJECT_DIR)/bin/oapi-codegen
OAPI_GEN_DIR=$(PROJECT_DIR)/internal/apigen
OAPI_CODEGEN_FIBER_BIN=$(PROJECT_DIR)/bin/oapi-codegen-fiber

install-oapi-codegen:
	DIR=$(PROJECT_DIR)/bin VERSION=${OAPI_CODEGEN_VERSION} ./scripts/install-oapi-codegen.sh
	
install-oapi-codegen-fiber:
	GOBIN=$(PROJECT_DIR)/bin go install github.com/cloudcarver/oapi-codegen-fiber@v0.5.1

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
	DIR=$(PROJECT_DIR)/bin VERSION=${WIRE_VERSION} ./scripts/install-wire.sh

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
	DIR=$(PROJECT_DIR)/bin VERSION=${SQLC_VERSION} ./scripts/install-sqlc.sh

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
	DIR=$(PROJECT_DIR)/bin VERSION=${MOCKGEN_VERSION} ./scripts/install-mockgen.sh

gen-mock: install-mockgen
	$(MOCKGEN_BIN) -source=internal/model/model.go -destination=internal/model/mock_gen.go -package=model

###################################################
### Common
###################################################

gen: gen-spec gen-querier gen-wire gen-mock gen-frontend-client
	@go mod tidy

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

build-web:
	cd web && pnpm run build

build-server:
	GOOS=linux GOARCH=amd64 go build -o ./bin/wavekit-server cmd/main.go

build-docker:
	docker build -f docker/Dockerfile.pgbundle -t cloudcarver/wavekit:v0.1.0-pgbundle .
	docker build -f docker/Dockerfile -t cloudcarver/wavekit:v0.1.0 .

docker-push:
	docker push cloudcarver/wavekit:v0.1.0-pgbundle
	docker push cloudcarver/wavekit:v0.1.0

build: build-web build-server build-docker

push: docker-push
	docker push cloudcarver/wavekit:v0.1.0-pgbundle

ut:
	COLOR=ALWAYS go test -race -covermode=atomic -coverprofile=coverage.out -tags ut ./... 
	@go tool cover -html coverage.out -o coverage.html
	@go tool cover -func coverage.out | fgrep total | awk '{print "Coverage:", $$3}'
