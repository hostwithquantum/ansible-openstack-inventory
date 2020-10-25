package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

// GetByCustomer ...
func GetByCustomer(provider *gophercloud.ProviderClient, customer string) []servers.Server {
	client, err := openstack.NewComputeV2(provider, gophercloud.EndpointOpts{Region: "RegionOne"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	allPages, err := servers.List(client, nil).AllPages()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var customerServers []servers.Server

	allServers, err := servers.ExtractServers(allPages)
	for _, server := range allServers {
		parts := strings.Split(server.Name, ".")
		if parts[len(parts)-1] == customer {
			customerServers = append(customerServers, server)
		}
	}

	return customerServers
}
