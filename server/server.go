package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
)

type ansibleServer struct {
	Name      string
	IPAddress string
}

// API ...
type API struct {
	accessNetwork string
	provider      *gophercloud.ProviderClient
}

// NewAPI ... factory/ctor
func NewAPI(network string, provider *gophercloud.ProviderClient) *API {
	api := new(API)
	api.accessNetwork = network
	api.provider = provider

	return api
}

// GetByCustomer ...
func (api API) GetByCustomer(customer string) []ansibleServer {
	// FIXME: if we use this multiple times, we should init it once
	client, err := openstack.NewComputeV2(api.provider, gophercloud.EndpointOpts{Region: "RegionOne"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	allPages, err := servers.List(client, nil).AllPages()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var customerServers []ansibleServer

	allServers, err := servers.ExtractServers(allPages)
	for _, server := range allServers {
		// server.CUSTOMER pattern
		parts := strings.Split(server.Name, ".")

		if parts[len(parts)-1] != customer {
			continue
		}

		node := ansibleServer{
			Name: server.Name,
		}

		for network, networkDetails := range server.Addresses {
			if network != api.accessNetwork {
				continue
			}

			for _, data := range networkDetails.([]interface{}) {
				for k, v := range data.(map[string]interface{}) {
					if k != "addr" {
						continue
					}

					node.IPAddress = v.(string)
					break
				}
			}
		}

		customerServers = append(customerServers, node)
	}

	return customerServers
}
