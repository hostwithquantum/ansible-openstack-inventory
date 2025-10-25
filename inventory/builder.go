package inventory

import (
	log "log/slog"
	"os"

	"github.com/hostwithquantum/ansible-openstack-inventory/host"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

func (inventory AnsibleInventory) BuildServers(nodes []server.AnsibleServer, defaultGroup string) {
	for _, node := range nodes {
		inventory.AddHostToGroup(node.Name, defaultGroup)

		hostVars := host.Build(node)
		for k, v := range hostVars {
			inventory.AddHostVar(k, v, node.Name)
		}
	}
}

func (inventory AnsibleInventory) BuildServerGroups(nodes []server.AnsibleServer, groups []string) {
	for _, node := range nodes {
		for _, g := range groups {
			switch g {
			case "all":
				log.Error("Group 'all' should not be part of the (server) groups.")
				os.Exit(1)

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

	// filter empty groups while we are at it
	for _, group := range inventory.Groups {
		// fmt.Printf("%s ------- %v", group.Name, group.Hosts)
		if len(group.Hosts) > 0 {
			jsonMap[group.Name] = inventory.Groups[group.Name]
		}
	}

	return jsonMap
}
