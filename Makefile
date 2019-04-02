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