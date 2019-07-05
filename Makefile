.PHONY: fmt golint test test-with-coverage ci
PACKAGES=`go list ./... | grep -v vendor | grep -v mocks`

fmt:
	for pkg in ${PACKAGES}; do \
		go fmt $$pkg; \
	done;

test:
	TEST_FAILED= ; \
	for pkg in ${PACKAGES}; do \
		go test $$pkg || TEST_FAILED=1; \
	done; \
	[ -z "$$TEST_FAILED" ]

test-with-coverage:
	echo "" > coverage.out
	echo "mode: set" > coverage-all.out
	TEST_FAILED= ; \
	for pkg in ${PACKAGES}; do \
		go test -coverprofile=coverage.out -covermode=set $$pkg || TEST_FAILED=1; \
		tail -n +2 coverage.out >> coverage-all.out; \
	done; \
	[ -z "$$TEST_FAILED" ]
	#go tool cover -html=coverage-all.out

ci:
	make test-with-coverage
