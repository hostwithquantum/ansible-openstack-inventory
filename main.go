package main

import (
	"encoding/json"
	"fmt"
	"os"

	"gopkg.in/ini.v1"

	"github.com/hostwithquantum/ansible-openstack-inventory/auth"
	"github.com/hostwithquantum/ansible-openstack-inventory/inventory"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

var defaultGroup = "all"

// FIXME: move into config.ini
var defaultVars = map[string]string{
	"ansible_ssh_user":           "core",
	"ansible_ssh_common_args":    "-F customers/files/quantum/ssh_config",
	"ansible_python_interpreter": "/opt/python/bin/python",
}

func main() {
	cfg, err := ini.Load("config.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	var accessNetwork = cfg.Section("").Key("network").String()
	var childrenGroups = cfg.Section("all").Key("children").Strings(",")

	provider, err := auth.Authenticate()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	customer, customerSet := os.LookupEnv("QUANTUM_CUSTOMER")
	if !customerSet {
		fmt.Println("Please set/export QUANTUM_CUSTOMER")
		os.Exit(1)
	}

	api := server.NewAPI(accessNetwork, provider)

	allServers := api.GetByCustomer(customer)

	inventory := inventory.NewInventory(customer, append(childrenGroups, defaultGroup))

	for defaultVar, defaultValue := range defaultVars {
		inventory.AddVarToGroup(defaultGroup, defaultVar, defaultValue)
	}

	for _, server := range allServers {
		inventory.AddHostToGroup(server.Name, defaultGroup)
		inventory.AddHostVar("ansible_host", server.IPAddress, server.Name)

		for _, g := range childrenGroups {
			inventory.AddHostToGroup(server.Name, g)
		}
	}

	for _, group := range append(childrenGroups, defaultGroup) {
		sec, err := cfg.GetSection(group)
		if err != nil {
			continue
		}

		inventory.AddChildrenToGroup(sec.Key("children").Strings(","), group)
	}

	// FIXME: move into config.ini
	inventory.AddVarToGroup("docker_swarm_manager", "swarm_labels", []string{
		"quantum",
		"manager",
		customer,
	})

	json, err := json.Marshal(inventory.ReturnJSONInventory())
	if err != nil {
		fmt.Println("Failed to encode")
	}

	fmt.Println(string(json))
	os.Exit(0)
}
