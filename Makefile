BUILD_DIR ?= bin
BUILD_PACKAGE ?= ./cmd/main.go
PROJECT_PKG = github.com/bogatyr285/auth-go

VERSION ?= $(shell git describe --tags --exact-match 2>/dev/null || git symbolic-ref -q --short HEAD)
COMMIT_HASH ?= $(shell git rev-parse --short HEAD 2>/dev/null)
BUILD_DATE ?= $(shell date +%FT%T%z)
# remove debug info from the binary & make it smaller
LDFLAGS += -s -w
LDFLAGS += -X ${PROJECT_PKG}/internal/buildinfo.version=${VERSION} -X ${PROJECT_PKG}/internal/buildinfo.commitHash=${COMMIT_HASH} -X ${PROJECT_PKG}/internal/buildinfo.buildDate=${BUILD_DATE}

build:
	go build ${GOARGS} -tags "${GOTAGS}" -ldflags "${LDFLAGS}" -o ${BUILD_DIR}/${BINARY_NAME} ${BUILD_PACKAGE}

proto:
	protoc \
	-I ./api/grpc \
	-I./third_party/googleapis \
	--go_out=./pkg/server/grpc/auth --go_opt=paths=source_relative \
	--go-grpc_out=./pkg/server/grpc/auth --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=./pkg/server/grpc/auth --grpc-gateway_opt=paths=source_relative,generate_unbound_methods=true \
	./api/grpc/*.proto


keys:
	ssh-keygen -t rsa -b 4096 -m PEM -f jwtRS256.key