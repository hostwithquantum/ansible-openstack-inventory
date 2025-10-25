package response_test

import (
	"errors"
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/response"
	"github.com/stretchr/testify/assert"
)

var (
	expectation = "{\"_meta\": {\"hostvars\": {}}}"
)

func TestError_BuildEmptyResponse(t *testing.T) {
	empty := response.BuildEmptyRepository(errors.New("foo"))

	// cannot return error in response -> Ansible barfs
	assert.Equal(t, expectation, empty, "Broken empty thing:\n %v\n (expected: %v)", empty, expectation)
}

func TestNoError_BuildEmptyResponse(t *testing.T) {
	empty := response.BuildEmptyRepository(nil)
	assert.Equal(t, expectation, empty, "Broken empty thing:\n %v\n (expected: %v)", empty, expectation)
}
