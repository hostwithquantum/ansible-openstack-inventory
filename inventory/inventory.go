package inventory

import (
	"encoding/json"
	"log"
)

// AnsibleInventory ...
type AnsibleInventory struct {
	hosts    []string
	groups   map[string]inventoryGroup
	hostvars map[string]map[string]string
}

type inventoryGroup struct {
	Name     string                 `json:"-"`
	Hosts    []string               `json:"hosts"`
	Vars     map[string]interface{} `json:"vars,omitempty"`
	Children []string               `json:"children,omitempty"`
}

// NewInventory ... factory/ctor
func NewInventory(customer string, groups []string) *AnsibleInventory {
	inventory := new(AnsibleInventory)

	inventory.hosts = make([]string, 0)
	inventory.groups = make(map[string]inventoryGroup)
	inventory.hostvars = make(map[string]map[string]string)

	for _, child := range groups {
		var group = inventoryGroup{
			Name: child,
		}

		group.Vars = make(map[string]interface{})

		inventory.groups[child] = group
	}

	return inventory
}

// AddChildrenToGroup ...
func (inventory AnsibleInventory) AddChildrenToGroup(groups []string, group string) {
	g := inventory.groups[group]
	g.Children = groups
	inventory.groups[group] = g
}

// AddHostToGroup ...
func (inventory AnsibleInventory) AddHostToGroup(host string, group string) {
	g := inventory.groups[group]
	g.Hosts = append(g.Hosts, host)
	inventory.groups[group] = g
}

// AddHostVar ...
func (inventory AnsibleInventory) AddHostVar(variable string, value string, host string) {
	_, ok := inventory.hostvars[host]
	if !ok {
		inventory.hostvars[host] = make(map[string]string)
	}

	inventory.hostvars[host][variable] = value
}

// AddVarToGroup ...
func (inventory AnsibleInventory) AddVarToGroup(group string, variable string, value interface{}) {
	g := inventory.groups[group]
	g.Vars[variable] = value
	inventory.groups[group] = g
}

// ReturnJSONInventory ...
func (inventory AnsibleInventory) ReturnJSONInventory() string {
	jsonMap := make(map[string]interface{})

	hostvars := make(map[string]map[string]map[string]string)
	hostvars["hostvars"] = inventory.hostvars

	jsonMap["_meta"] = hostvars
	jsonMap["all"] = inventory.groups["all"]

	for _, group := range inventory.groups {
		jsonMap[group.Name] = inventory.groups[group.Name]
	}

	jsonByte, err := json.Marshal(jsonMap)
	if err != nil {
		log.Fatal(err)
	}

	return string(jsonByte)
}
