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
		fmt.Println(server.Name)
		for network, networkDetails := range server.Addresses {
			if network == "quantum-internal" {
				fmt.Printf("%v\n", networkDetails)
			}
		}
	}
}
