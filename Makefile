.PHONY: build
build:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: dev
dev:
	go run main.go | jq '.'

.PHONY: run
run: build
	./dist/ansible-openstack-inventory_darwin_amd64/ansible-openstack-inventory --list | jq '.'

test:
	go test ./...
