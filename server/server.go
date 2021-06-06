package server

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/pagination"
)

// AnsibleServer ... simple struct to build an inventory from
type AnsibleServer struct {
	Name       string
	IPAddress  string
	FloatingIP string
	MetaData   map[string]string
}

// API ...
type API struct {
	accessNetwork string
	provider      *gophercloud.ProviderClient
	client        *gophercloud.ServiceClient
}

// NewAPI ... factory/ctor
func NewAPI(network string, provider *gophercloud.ProviderClient) *API {
	api := new(API)
	api.accessNetwork = network
	api.provider = provider

	client, err := openstack.NewComputeV2(api.provider, gophercloud.EndpointOpts{Region: "RegionOne"})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api.client = client

	return api
}

// GetByNode ...
func (api API) GetByNode(host string) (AnsibleServer, error) {
	opts := servers.ListOpts{Name: host}
	pager := servers.List(api.client, opts)

	server := AnsibleServer{}

	publicIPs := api.getFloatingIps()

	err := pager.EachPage(func(page pagination.Page) (bool, error) {
		serverList, err := servers.ExtractServers(page)
		if err != nil {
			return false, err
		}

		if len(serverList) == 0 {
			return false, nil
		}

		for _, s := range serverList {
			server.Name = s.Name
			server.IPAddress = extractIP(s.Addresses, api.accessNetwork)
			server.MetaData = s.Metadata

			if _, ok := publicIPs[s.ID]; ok {
				server.FloatingIP = publicIPs[s.ID]
			}
		}

		return false, nil
	})

	if err != nil {
		return server, err
	}

	// empty result
	if server.Name == "" {
		return server, fmt.Errorf("Could not find a host named: %s", host)
	}

	return server, nil
}

// GetByCustomer ...
func (api API) GetByCustomer(customer string) []AnsibleServer {
	publicIPs := api.getFloatingIps()

	allPages, err := servers.List(api.client, nil).AllPages()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	var customerServers []AnsibleServer

	allServers, err := servers.ExtractServers(allPages)
	for _, server := range allServers {
		// server.CUSTOMER pattern
		parts := strings.Split(server.Name, ".")

		if parts[len(parts)-1] != customer {
			continue
		}

		node := AnsibleServer{
			Name:      server.Name,
			IPAddress: extractIP(server.Addresses, api.accessNetwork),
			MetaData:  server.Metadata,
		}

		if _, ok := publicIPs[server.ID]; ok {
			node.FloatingIP = publicIPs[server.ID]
		}

		customerServers = append(customerServers, node)
	}

	return customerServers
}

func (api API) getFloatingIps() map[string]string {
	allPages, err := floatingips.List(api.client).AllPages()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	allFloatingIPs, err := floatingips.ExtractFloatingIPs(allPages)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	keep := make(map[string]string)
	for _, floatingIP := range allFloatingIPs {
		keep[floatingIP.InstanceID] = floatingIP.IP
	}

	return keep
}

func extractIP(addresses map[string]interface{}, network string) string {
	for networkName, networkDetails := range addresses {
		if networkName != network {
			continue
		}

		for _, data := range networkDetails.([]interface{}) {
			for k, v := range data.(map[string]interface{}) {
				if k != "addr" {
					continue
				}

				return v.(string)
			}
		}
	}

	return ""
}
