OUT_DIR=bin
build:
	mkdir -p "${OUT_DIR}"
	export GO111MODULE="on" && export GOFLAGS="" && export GOWORK=off && go build  -ldflags="-s -w" -mod=mod -o "${OUT_DIR}" .
run: build
	./bin/test-utils
