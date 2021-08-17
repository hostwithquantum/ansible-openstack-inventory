package lbaas

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
)

type API struct {
	customer string
	client   *gophercloud.ServiceClient
	provider *gophercloud.ProviderClient
}

const OS_LB_PROVIDER string = "octavica"

func NewAPI(customer string, provider *gophercloud.ProviderClient) *API {
	api := new(API)
	api.provider = provider
	api.customer = customer

	client, err := openstack.NewLoadBalancerV2(provider, gophercloud.EndpointOpts{
		Name:   "octavia",
		Region: "RegionOne",
	})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	api.client = client

	return api
}

func (api API) GetAll() []loadbalancers.LoadBalancer {
	allPages, err := loadbalancers.List(api.client, nil).AllPages()
	if err != nil {
		fmt.Printf("First: %s", err)
		os.Exit(1)
	}
	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		fmt.Printf("Second: %s", err)
		os.Exit(1)
	}

	return allLoadbalancers
}

func (api API) GetById(id string) (loadbalancers.LoadBalancer, error) {
	listOpts := loadbalancers.ListOpts{
		ID: id,
	}

	return api.doSingleRequest(listOpts)
}

func (api API) GetByName() (loadbalancers.LoadBalancer, error) {
	listOpts := loadbalancers.ListOpts{
		Name: buildLBName(api.customer),
	}

	return api.doSingleRequest(listOpts)
}

func (api API) HasLB() bool {
	listOpts := loadbalancers.ListOpts{
		Name: buildLBName(api.customer),
	}

	_, err := api.doSingleRequest(listOpts)
	if err != nil {
		return false
	}

	return true
}

func (api API) doSingleRequest(listOpts loadbalancers.ListOpts) (loadbalancers.LoadBalancer, error) {
	var lb loadbalancers.LoadBalancer
	allPages, err := loadbalancers.List(api.client, listOpts).AllPages()
	if err != nil {
		return lb, err
	}

	allLoadbalancers, err := loadbalancers.ExtractLoadBalancers(allPages)
	if err != nil {
		return lb, err
	}

	if len(allLoadbalancers) == 0 {
		return lb, fmt.Errorf("Couldn't find loadbalancer: %v", listOpts)
	}

	return allLoadbalancers[0], nil
}

// The pattern matches that of:
// https://github.com/hostwithquantum/terraform-openstack-loadbalancer/blob/dev-clean-up/main.tf#L5
func buildLBName(customer string) string {
	return fmt.Sprintf("%s-loadbalancer", customer)
}
