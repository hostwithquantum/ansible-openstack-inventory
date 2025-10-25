package utils_test

import (
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/internal/utils"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
	"github.com/stretchr/testify/assert"
)

func TestSortServersByName_Sequential(t *testing.T) {
	servers := []server.AnsibleServer{
		{Name: "node-003.cluster"},
		{Name: "node-001.cluster"},
		{Name: "node-002.cluster"},
	}

	utils.SortServersByName(servers)

	expected := []string{"node-001.cluster", "node-002.cluster", "node-003.cluster"}
	actual := make([]string, len(servers))
	for i, s := range servers {
		actual[i] = s.Name
	}

	assert.Equal(t, expected, actual, "Servers should be sorted in ascending order")
}

func TestSortServersByName_ReverseOrder(t *testing.T) {
	servers := []server.AnsibleServer{
		{Name: "node-006.cluster"},
		{Name: "node-005.cluster"},
		{Name: "node-004.cluster"},
	}

	utils.SortServersByName(servers)

	expected := []string{"node-004.cluster", "node-005.cluster", "node-006.cluster"}
	actual := make([]string, len(servers))
	for i, s := range servers {
		actual[i] = s.Name
	}

	assert.Equal(t, expected, actual, "Servers should be sorted in ascending order")
}

func TestSortServersByName_MixedOrder(t *testing.T) {
	servers := []server.AnsibleServer{
		{Name: "node-005.cluster"},
		{Name: "node-001.cluster"},
		{Name: "node-006.cluster"},
		{Name: "node-002.cluster"},
		{Name: "node-004.cluster"},
		{Name: "node-003.cluster"},
	}

	utils.SortServersByName(servers)

	expected := []string{
		"node-001.cluster",
		"node-002.cluster",
		"node-003.cluster",
		"node-004.cluster",
		"node-005.cluster",
		"node-006.cluster",
	}
	actual := make([]string, len(servers))
	for i, s := range servers {
		actual[i] = s.Name
	}

	assert.Equal(t, expected, actual, "Servers should be sorted in ascending order")
}
