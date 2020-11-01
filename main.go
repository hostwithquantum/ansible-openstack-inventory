package main

import (
	"fmt"
	"os"

	"github.com/hostwithquantum/ansible-openstack-inventory/auth"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

func main() {
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

	allServers := server.GetByCustomer(provider, customer)
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
	}
}
