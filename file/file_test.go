package file_test

import (
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/file"
	"github.com/stretchr/testify/assert"
)

func Test_HandleGroup(t *testing.T) {
	gvf := file.NewGroupVarsFile("./data")
	groupFileYaml, err := gvf.HandleGroup("group")
	assert.NoError(t, err)

	tests := make(map[string]interface{})
	tests["variable"] = "value"
	tests["lala"] = 1.0

	for varKey, varValue := range tests {
		val, ok := groupFileYaml[varKey]
		assert.True(t, ok, "Could not find '%s'", varKey)
		assert.Equal(t, varValue, val, "Value doesn't match: %v (expected: %v)", val, varValue)
	}
}
