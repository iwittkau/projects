.PHONY: run-dev
run-dev:
	gow run cmd/projects/main.go

.PHONY: run-dev-git
run-dev-git:
	go run cmd/projects/main.go --git

install:
	go install github.com/iwittkau/projects/cmd/projects