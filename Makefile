bin:=./dist/ansible-openstack-inventory_darwin_amd64/ansible-openstack-inventory


.PHONY: build
build:
	goreleaser --snapshot --skip-publish --rm-dist

.PHONY: dev
dev:
	go run main.go | jq '.'

.PHONY: e2e
e2e: build
	$(bin) --list | jq -r '.' > e2e/result.json

.PHONY: run
run: build
	$(bin) --list | jq '.'

test:
	go test ./...
