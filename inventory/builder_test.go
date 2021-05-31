package inventory_test

import (
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/inventory"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

func Test_Cluster(t *testing.T) {
	servers := []server.AnsibleServer{
		createManager(),
		createWorker(),
	}

	groups := []string{
		"docker_swarm_manager",
		"docker_swarm_worker",
		"portainer-agent-node",
		"promtail",
	}

	ansible := inventory.NewInventory("customer_name", append(groups, "all"))
	ansible.BuildServers(servers, "all")
	ansible.BuildServerGroups(servers, groups)

	jsonInventory := ansible.BuildInventory()
	// fmt.Printf("%v", jsonInventory)
	// fmt.Printf("%v", jsonInventory["_meta"])

	manager_group := jsonInventory["docker_swarm_manager"].(inventory.InventoryGroup)
	worker_group := jsonInventory["docker_swarm_worker"].(inventory.InventoryGroup)

	assertHostCount(manager_group, 1, t)
	assertHostCount(worker_group, 1, t)
	assertHostCount(jsonInventory["all"].(inventory.InventoryGroup), 2, t)

	host_data := jsonInventory["_meta"]
	// fmt.Printf("%v", host_data)

	for _, s := range servers {
		assertHostVarsLabels(host_data, s.Name, "foo", t)
	}

	for _, g := range groups {
		_g, ok := jsonInventory[g].(inventory.InventoryGroup)
		if !ok {
			t.Errorf("Could not find group: %s", _g)
		}

		switch g {
		case "docker_swarm_manager":
		case "docker_swarm_worker":
		case "all":
			break

		default:
			assertHostCount(_g, len(servers), t)
		}
	}

	assertGroupLabels(manager_group, "manager", t)
	assertGroupLabels(worker_group, "worker", t)
	assertGroupLabels(manager_group, "quantum", t)
	assertGroupLabels(worker_group, "quantum", t)

}

func assertHostCount(group inventory.InventoryGroup, count int, t *testing.T) {
	if len(group.Hosts) != count {
		t.Errorf("%s needs to have %d host(s)", group.Name, count)
	}
}

func assertGroupLabels(group inventory.InventoryGroup, label string, t *testing.T) {
	all_labels, ok := group.Vars["swarm_labels"].([]string)
	if !ok {
		t.Errorf("No labels found for group %s", group.Name)
	}

	if len(all_labels) == 0 {
		t.Errorf("Labels are empty for group %s", group.Name)
	}

	for _, l := range all_labels {
		if l == label {
			return
		}
	}

	t.Errorf("Found labels, but not the label (%s) we needed for group %s", label, group.Name)
}

func assertHostVarsLabels(meta interface{}, node, label string, t *testing.T) {
	// This is definitely an obscure test
	_label, ok := meta.(map[string]map[string]map[string]string)["hostvars"][node]["swarm_labels"]
	if !ok {
		t.Errorf("Could not find swarm_labels for node %s", node)
	}
	if label == _label {
		return
	}

	t.Errorf("Could not find label %s on node %s", label, node)
}

func createManager() server.AnsibleServer {
	return server.AnsibleServer{
		Name:       "node-001.customer",
		IPAddress:  "192.168.1.10",
		FloatingIP: "10.10.10.10",
		MetaData: map[string]string{
			"com.planetary-quantum.meta.role":  "manager",
			"com.planetary-quantum.meta.label": "foo",
		},
	}
}

func createWorker() server.AnsibleServer {
	return server.AnsibleServer{
		Name:      "node-002.customer",
		IPAddress: "192.168.1.11",
		MetaData: map[string]string{
			"com.planetary-quantum.meta.role":  "worker",
			"com.planetary-quantum.meta.label": "foo",
		},
	}
}
