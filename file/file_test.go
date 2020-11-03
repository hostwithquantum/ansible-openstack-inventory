package file_test

import (
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/file"
)

func Test_HandleGroup(t *testing.T) {
	gvf := file.NewGroupVarsFile("./data")
	groupFileYaml, err := gvf.HandleGroup("group")
	if err != nil {
		t.Error(err)
	}

	tests := make(map[string]interface{})
	tests["variable"] = "value"
	tests["lala"] = 1.0

	for varKey, varValue := range tests {
		val, ok := groupFileYaml[varKey]
		if !ok {
			t.Errorf("Could not find '%s'", varKey)
		}

		if varValue != val {
			t.Errorf("Value doesn't match: %v (expected: %v)", val, varValue)
		}
	}
}
