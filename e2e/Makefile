inv_build:=../dist/ansible-openstack-inventory_darwin_amd64/ansible-openstack-inventory

.PHONY: e2e
e2e:
	$(info You can run `make test-host` or `make test-no-host`)

.PHONY: _dep
_dep:
	cd ../ && make build

.PHONY: all
all: _dep test-host test-no-host

.PHONY: test-host
test-host:
	$(info This should match hosts)
	ansible \
		-i $(inv_build) \
		-m ping all || true

.PHONY: test-no-host
test-no-host:
	$(info This shouldn't return any)
	OS_TENANT_NAME=${EMPTY_TENANT} ansible \
		-i $(inv_build) \
		-m ping all || true
