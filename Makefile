.PHONY: build run

build:
	goreleaser --snapshot --skip-publish --rm-dist

run: build
	./ansible-openstack-inventory | jq '.'
