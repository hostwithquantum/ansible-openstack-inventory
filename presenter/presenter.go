package presenter

import (
	"strings"

	"github.com/gophercloud/gophercloud/openstack/loadbalancer/v2/loadbalancers"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

const role_label string = "com.planetary-quantum.meta.role"

type Presenter struct {
	FIPs map[string]map[string]string
}

func (p Presenter) AddFipToNode(node server.AnsibleServer) server.AnsibleServer {
	if _, ok := p.FIPs[node.ID]; ok {
		node.FloatingIP = p.FIPs[node.ID]["ip"]
	}

	return node
}

func (p Presenter) AddFipsToNodes(nodes []server.AnsibleServer) []server.AnsibleServer {
	for idx, node := range nodes {
		nodes[idx] = p.AddFipToNode(node)
	}

	return nodes
}

func (p Presenter) AddLoadbalancerToNode(node server.AnsibleServer, lb loadbalancers.LoadBalancer) server.AnsibleServer {
	// is this a manager?
	if isFirstManagerWithoutFip(node) {
		// fmt.Println("NO FIP SO FAR")
		// fmt.Printf("LB ID: %s\n", lb.ID)
		// fmt.Printf("%v\n", p.FIPs)
		for fipInstanceID, fip := range p.FIPs {
			if fip["target"] == "lb" && lb.ID == fipInstanceID {
				node.FloatingIP = fip["ip"]
			}
		}
	}

	return node
}

func isFirstManagerWithoutFip(node server.AnsibleServer) bool {
	if !strings.Contains(node.Name, "node-001") {
		return false
	}

	if node.MetaData[role_label] != "manager" {
		return false
	}

	if node.FloatingIP == "" {
		return true
	}

	return false
}
