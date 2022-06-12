TEST?=$$(go list ./... |grep -v 'vendor')
GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)

default: build

build: fmt
	./scripts/build.sh

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

test: fmt
	go test $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

acceptance: fmt
	go test -v $(TEST) || exit 1
	echo $(TEST) | \
		TF_ACC=1 xargs -t -n4 go test -v $(TESTARGS) -parallel=4

start-fakeserver:
	@sh -c "$(CURDIR)/scripts/fakeserver.sh start"

stop-fakeserver:
	@sh -c "$(CURDIR)/scripts/fakeserver.sh stop"
