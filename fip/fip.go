package fip

import (
	"fmt"
	"os"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/floatingips"
)

type FIP struct {
	accessNetwork string
	provider      *gophercloud.ProviderClient
	client        *gophercloud.ServiceClient
}

func NewFIP(network string, provider *gophercloud.ProviderClient) *FIP {
	api := new(FIP)
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

func (api FIP) GetIps() map[string]map[string]string {
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

	keep := make(map[string]map[string]string)
	for _, floatingIP := range allFloatingIPs {
		//fmt.Printf("%v", floatingIP)

		id := strings.Replace(floatingIP.InstanceID, "lb-", "", 1)

		keep[id] = make(map[string]string)

		keep[id]["ip"] = floatingIP.IP
		if strings.Contains(floatingIP.InstanceID, "lb-") {
			keep[id]["target"] = "lb"
		} else {
			keep[id]["target"] = "vm"
		}
	}

	//fmt.Printf("%v", keep)

	return keep
}
