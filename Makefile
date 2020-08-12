# Borrowed from:
# https://gist.github.com/turtlemonvh/38bd3d73e61769767c35931d8c70ccb4
# https://github.com/silven/go-example/blob/master/Makefile
# https://vic.demuzere.be/articles/golang-makefile-crosscompile/

BINARY=den-bot
GOARCH=amd64

CURRENT_DIR=$(shell pwd)

COMMIT=$(shell git rev-parse HEAD)
BRANCH=$(shell git rev-parse --abbrev-ref HEAD)
VERSION=$(shell cat .version)
METRICS_IMPORT_PATH=github.com/caquillo07/rotom-bot/metrics

# Setup the -ldflags option for go build here, interpolate the variable values
LDFLAGS = -ldflags "-X ${METRICS_IMPORT_PATH}.Version=${VERSION} -X ${METRICS_IMPORT_PATH}.Commit=${COMMIT} -X ${METRICS_IMPORT_PATH}.Branch=${BRANCH}"

dev:
	go run main.go bot --dev-log

dev-reload:
	air -c .air.conf

migrate-dev:
	go run main.go migrate --config example-config.yaml --dev-log

linux:
	GOOS=linux GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-linux-${GOARCH} .

darwin:
	GOOS=darwin GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-darwin-${GOARCH} .

windows:
	GOOS=windows GOARCH=${GOARCH} go build ${LDFLAGS} -o ${BINARY}-windows-${GOARCH} .

release:
	git tag -a ${VERSION} -m "release-${VERSION}"
	git push --force origin ${VERSION}
	goreleaser --rm-dist