.PHONY: build run

build:
	go build

run: build
	./ansible-openstack-inventory | jq '.'
