package host

import (
	log "github.com/sirupsen/logrus"

	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

func Build(node server.AnsibleServer) map[string]string {
	hostVars := make(map[string]string)
	hostVars["ansible_host"] = node.IPAddress

	if node.FloatingIP != "" {
		hostVars["floating_ip"] = node.FloatingIP
	}

	if len(node.MetaData) == 0 {
		log.Debugf("Empty/broken metadata for %s", node.Name)
		return hostVars
	}

	if node_swarm_label, ok := node.MetaData["com.planetary-quantum.meta.label"]; ok {
		hostVars["swarm_labels"] = node_swarm_label
	}

	node_group, err := server.GetGroup(node)
	if err != nil {
		log.Debug(err)
	} else {
		hostVars["quantum_group_name"] = node_group
	}

	return hostVars
}
