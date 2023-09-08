all: install

install: go.sum
	#go install -tags "${build_tags}" ./cmd/nsd
	go install ./cmd/nsd
	#go install -tags "${build_tags}" ./cmd/nscli
	go install ./cmd/nscli

go.sum: go.mod
	@echo "--> Ensure dependencies have not been modified"
	@go mod verify

go.mod:
	go mod tidy