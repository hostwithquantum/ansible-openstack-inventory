package inventory_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

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

func createBasicInventory() *inventory.AnsibleInventory {
	inventory := inventory.NewInventory("customer", []string{"all", "test"})
	inventory.AddHostToGroup("127.0.0.1", "all")
	inventory.AddHostToGroup("127.0.0.1", "test")
	inventory.AddHostVar("host_var", "value", "127.0.0.1")
	inventory.AddVarToGroup("test", "group_var", []string{"bar", "foobar"})

	return inventory
}
