package inventory

import (
	"encoding/json"
	"log"
)

// AnsibleInventory ...
type AnsibleInventory struct {
	Hosts    []string
	Groups   map[string]InventoryGroup
	Hostvars map[string]map[string]string
}

// InventoryGroup ...
type InventoryGroup struct {
	Name     string                 `json:"-"`
	Hosts    []string               `json:"hosts"`
	Vars     map[string]interface{} `json:"vars,omitempty"`
	Children []string               `json:"children,omitempty"`
}

// NewInventory ... factory/ctor
func NewInventory(customer string, groups []string) *AnsibleInventory {
	inventory := new(AnsibleInventory)

	inventory.Hosts = make([]string, 0)
	inventory.Groups = make(map[string]InventoryGroup)
	inventory.Hostvars = make(map[string]map[string]string)

	for _, child := range groups {
		var group = InventoryGroup{
			Name: child,
		}

		group.Vars = make(map[string]interface{})

		inventory.Groups[child] = group
	}

	return inventory
}

// AddChildrenToGroup ...
func (inventory AnsibleInventory) AddChildrenToGroup(groups []string, group string) {
	g := inventory.Groups[group]
	g.Children = groups
	inventory.Groups[group] = g
}

// AddHostToGroup ...
func (inventory AnsibleInventory) AddHostToGroup(host string, group string) {
	g := inventory.Groups[group]
	g.Hosts = append(g.Hosts, host)
	inventory.Groups[group] = g
}

// AddHostVar ...
func (inventory AnsibleInventory) AddHostVar(variable string, value string, host string) {
	_, ok := inventory.Hostvars[host]
	if !ok {
		inventory.Hostvars[host] = make(map[string]string)
	}

	inventory.Hostvars[host][variable] = value
}

// AddVarToGroup ...
func (inventory AnsibleInventory) AddVarToGroup(group string, variable string, value interface{}) {
	g := inventory.Groups[group]
	g.Vars[variable] = value
	inventory.Groups[group] = g
}

// ReturnJSONInventory ...
func (inventory AnsibleInventory) ReturnJSONInventory() string {
	jsonMap := inventory.BuildInventory()

	jsonByte, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonByte)
}
