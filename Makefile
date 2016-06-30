
TARGETS_NOVENDOR := $(shell glide novendor)

fmt:
	@echo $(TARGETS_NOVENDOR) | xargs go fmt
