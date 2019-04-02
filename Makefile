APP_NAME = r3x-cli
ORG_NAME = rubixFunctions
PKG = github.com/$(ORG_NAME)/$(APP_NAME)
TOP_SRC_DIRS = cmd
PACKAGES     ?= $(shell sh -c "find $(TOP_SRC_DIRS) -name \\*_test.go \
                   -exec dirname {} \\; | sort | uniq")


.PHONY: test-all
test-all: test-unit
	make test-integration

.PHONY: test
test: test-unit

.PHONY: test-unit
test-unit:
	@echo Running tests:
	GOCACHE=off go test -cover \
	  $(addprefix $(PKG)/,$(PACKAGES))

.PHONY: test-integration
test-integration:
	@echo Running tests:
	GOCACHE=off go test -failfast -cover -tags=integration \
	  $(addprefix $(PKG)/,$(PACKAGES))

.PHONY: test-integration-cover
test-integration-cover:
	echo "mode: count" > coverage-all.out
	GOCACHE=off $(foreach pkg,$(PACKAGES),\
		go test -failfast -tags=integration -coverprofile=coverage.out -covermode=count $(addprefix $(PKG)/,$(pkg)) || exit 1;\
		tail -n +2 coverage.out >> coverage-all.out;)
	make cleanup-coverage-file

.PHONY: cleanup-coverage-file
cleanup-coverage-file:
	@echo "Cleaning up output of coverage report"
	./scripts/cleanup-coverage-file.sh