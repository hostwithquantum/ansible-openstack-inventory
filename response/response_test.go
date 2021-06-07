package response_test

import (
	"errors"
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/response"
)

func TestError_BuildEmptyResponse(t *testing.T) {
	err := errors.New("foo")

	empty := response.BuildEmptyRepository(err)

	// cannot return error in response -> Ansible barfs
	expectation := "{\"_meta\": {\"hostvars\": {}}}"

	if empty != expectation {
		t.Errorf("Broken empty thing:\n %v\n (expected: %v)", empty, expectation)
	}
}

func TestNoError_BuildEmptyResponse(t *testing.T) {
	empty := response.BuildEmptyRepository(nil)

	expectation := "{\"_meta\": {\"hostvars\": {}}}"

	if empty != expectation {
		t.Errorf("Broken empty thing:\n %v\n (expected: %v)", empty, expectation)
	}
}
