
TARGETS_NOVENDOR := $(shell glide novendor)

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt

test:
	go test -v ./

coverage:
	go test -v -cover -covermode=count -coverprofile=coverage.txt ./

html: coverage
	go tool cover -html=coverage.txt && unlink coverage.txt
