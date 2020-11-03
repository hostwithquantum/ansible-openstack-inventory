.PHONY: build run

build:
	goreleaser --snapshot --skip-publish --rm-dist

run: build
	./dist/ansible-openstack-inventory_darwin_amd64/ansible-openstack-inventory --list | jq '.'

test:
	go test ./...
