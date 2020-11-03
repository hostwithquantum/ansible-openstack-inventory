package inventory_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/file"
	"github.com/hostwithquantum/ansible-openstack-inventory/inventory"
)

func ExampleNewInventory() {
	inventory := createBasicInventory()
	fmt.Println(inventory)
	// Output: &{[] map[all:{all [127.0.0.1] map[] []} test:{test [127.0.0.1] map[group_var:[bar foobar]] []}] map[127.0.0.1:map[host_var:value]]}
}

func Test_ReturnJSONInventory(t *testing.T) {
	inventory := createBasicInventory()
	json := inventory.ReturnJSONInventory()

	tests := []string{
		"{\"_meta\":",
		"{\"hostvars\":{\"127.0.0.1\":{\"host_var\":\"value\"}}}",
		"\"all\":{\"hosts\":[\"127.0.0.1\"]}",
		"\"test\":{\"hosts\":[\"127.0.0.1\"],\"vars\":{\"group_var\":[\"bar\",\"foobar\"]}}",
	}

	for _, test := range tests {
		ok := strings.Contains(json, test)
		if !ok {
			t.Log(json)
			t.Errorf("The generated JSON doesn't contain '%s'", test)
		}
	}
}

func Test_AddVarToGroup(t *testing.T) {
	inventory := createBasicInventory()

	groupName := "test"
	t.Logf("Testing '%s'", groupName)

	gvf := file.NewGroupVarsFile("./data")
	groupFileYaml, err := gvf.HandleGroup(groupName)
	if err != nil {
		t.Error(err)
	}

	for varKey, varValue := range groupFileYaml {
		inventory.AddVarToGroup(groupName, varKey, varValue)
	}

	json := inventory.ReturnJSONInventory()

	tests := []string{
		"\"and_another\":[\"this\",\"is\",\"a\",\"slice\"]",
		"\"another\":1",
		"\"some_variable\":\"string\"",
	}

	for _, test := range tests {
		ok := strings.Contains(json, test)
		if !ok {
			t.Errorf("JSON does not contain: '%s'", test)
		}
	}
}

func createBasicInventory() *inventory.AnsibleInventory {
	inventory := inventory.NewInventory("customer", []string{"all", "test"})
	inventory.AddHostToGroup("127.0.0.1", "all")
	inventory.AddHostToGroup("127.0.0.1", "test")
	inventory.AddHostVar("host_var", "value", "127.0.0.1")
	inventory.AddVarToGroup("test", "group_var", []string{"bar", "foobar"})

	return inventory
}
