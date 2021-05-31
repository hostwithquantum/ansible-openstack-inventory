package inventory

import (
	"log"

	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

func (inventory AnsibleInventory) BuildServers(nodes []server.AnsibleServer, defaultGroup string) {
	for _, node := range nodes {
		inventory.AddHostToGroup(node.Name, defaultGroup)
		inventory.AddHostVar("ansible_host", node.IPAddress, node.Name)

		if node.FloatingIP != "" {
			inventory.AddHostVar("floating_ip", node.FloatingIP, node.Name)
		}

		if len(node.MetaData) > 0 {
			node_swarm_label, ok := node.MetaData["com.planetary-quantum.meta.label"]
			if ok {
				inventory.AddHostVar("swarm_labels", node_swarm_label, node.Name)
			}
		}
	}
}

func (inventory AnsibleInventory) BuildServerGroups(nodes []server.AnsibleServer, groups []string) {
	for _, node := range nodes {
		for _, g := range groups {
			switch g {
			case "all":
				log.Fatal("Group 'all' should not be part of the (server) groups.")

			case "docker_swarm_manager":
				if server.IsManager(node) {
					labels := []string{
						"quantum",
						"manager",
					}

					inventory.AddHostToGroup(node.Name, g)
					inventory.AddVarToGroup(g, "swarm_labels", labels)
				}
				break
			case "docker_swarm_worker":
				if server.IsWorker(node) {
					labels := []string{
						"quantum",
						"worker",
					}
					inventory.AddHostToGroup(node.Name, g)
					inventory.AddVarToGroup(g, "swarm_labels", labels)
				}
				break
			default:
				inventory.AddHostToGroup(node.Name, g)
			}
		}
	}
}

func (inventory AnsibleInventory) BuildInventory() map[string]interface{} {
	jsonMap := make(map[string]interface{})

	hostvars := make(map[string]map[string]map[string]string)
	hostvars["hostvars"] = inventory.Hostvars

	jsonMap["_meta"] = hostvars
	jsonMap["all"] = inventory.Groups["all"]

	for _, group := range inventory.Groups {
		jsonMap[group.Name] = inventory.Groups[group.Name]
	}

	return jsonMap
}
