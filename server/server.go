package server

import (
	"fmt"
	"os"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
)

// AnsibleServer ... simple struct to build an inventory from
type AnsibleServer struct {
	ID         string // OpenStack Instance ID
	Name       string // OpenStack Instance Name
	IPAddress  string // IP address on customer network
	FloatingIP string // Floating IP
	MetaData   map[string]string
}

// API ...
type API struct {
	accessNetwork string
	customer      string
	provider      *gophercloud.ProviderClient
	client        *gophercloud.ServiceClient
	publicIPs     map[string]map[string]string
	lb            loadbalancers.LoadBalancer
}

// NewAPI ... factory/ctor
func NewAPI(customer string, network string, provider *gophercloud.ProviderClient) *API {
	api := new(API)
	api.accessNetwork = network
	api.provider = provider
	api.customer = customer

	client, err := openstack.NewComputeV2(api.provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api.client = client

	return api
}

// GetByNode ...
func (api API) GetByNode(host string) (AnsibleServer, error) {
	listOpts := servers.ListOpts{Name: host}
	servers, err := api.doRequest(listOpts)
	if err != nil {
		return AnsibleServer{}, err
	}

	if len(servers) == 0 {
		return AnsibleServer{}, fmt.Errorf("could not find a host named: %s", host)
	}

	return servers[0], nil
}

// GetByCustomer ...
func (api API) GetByCustomer(customer string) ([]AnsibleServer, error) {
	listOpts := servers.ListOpts{}
	return api.doRequest(listOpts)
}

func (api API) doRequest(listOpts servers.ListOpts) ([]AnsibleServer, error) {
	var customerServers []AnsibleServer

	allPages, err := servers.List(api.client, nil).AllPages()
	if err != nil {
		return customerServers, err
	}

	allServers, err := servers.ExtractServers(allPages)
	if err != nil {
		return customerServers, err
	}

	if len(allServers) == 0 {
		log.Debug("Couldn't find any servers in the tenant.")
	}

	for _, server := range allServers {
		// server.CUSTOMER pattern
		parts := strings.Split(server.Name, ".")

		if parts[len(parts)-1] != api.customer {
			continue
		}

		node := AnsibleServer{
			ID:        server.ID,
			Name:      server.Name,
			IPAddress: extractIP(server.Addresses, api.accessNetwork),
			MetaData:  server.Metadata,
		}

		customerServers = append(customerServers, node)
	}

	return customerServers, nil
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
