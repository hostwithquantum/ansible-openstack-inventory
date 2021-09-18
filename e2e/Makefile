inv_build:=../dist/ansible-openstack-inventory_darwin_amd64/ansible-openstack-inventory


assert.sh:
	wget \
		-O ./assert.sh \
		https://raw.githubusercontent.com/torokmark/assert.sh/main/assert.sh

.PHONY: clean
clean:
	rm assert.sh


.PHONY: e2e
e2e: assert.sh
	./basic-test.sh

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
