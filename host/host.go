package host

import (
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

func Build(node server.AnsibleServer) map[string]string {
	hostVars := make(map[string]string)
	hostVars["ansible_host"] = node.IPAddress

	if node.FloatingIP != "" {
		hostVars["floating_ip"] = node.FloatingIP
	}

	if len(node.MetaData) > 0 {
		node_swarm_label, ok := node.MetaData["com.planetary-quantum.meta.label"]
		if ok {
			hostVars["swarm_labels"] = node_swarm_label
		}
	}

	return hostVars
}
