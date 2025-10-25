package host_test

import (
	"reflect"
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/host"
	"github.com/hostwithquantum/ansible-openstack-inventory/server"
	"github.com/stretchr/testify/assert"
)

type testFixture struct {
	Node     server.AnsibleServer
	HostVars map[string]string
}

func Test_Build(t *testing.T) {
	data := testBuildDataprovider()
	for _, fixture := range data {
		hv := host.Build(fixture.Node)
		assert.True(t, reflect.DeepEqual(fixture.HostVars, hv), "Build is broken, expected '%v', but got '%v'", fixture.HostVars, hv)
	}
}

func testBuildDataprovider() []testFixture {
	// fixture 1
	server1 := server.AnsibleServer{
		IPAddress: "127.0.0.1",
		MetaData: map[string]string{
			"com.planetary-quantum.meta.label": "test",
		},
	}

	fixture1 := make(map[string]string)
	fixture1["ansible_host"] = "127.0.0.1"
	fixture1["swarm_labels"] = "test"

	server2 := server.AnsibleServer{
		IPAddress:  "127.0.0.2",
		FloatingIP: "192.168.1.1",
	}
	fixture2 := make(map[string]string)
	fixture2["ansible_host"] = "127.0.0.2"
	fixture2["floating_ip"] = "192.168.1.1"

	server3 := server.AnsibleServer{
		IPAddress: "127.0.0.3",
		MetaData: map[string]string{
			"com.planetary-quantum.meta.label": "google-dns",
		},
		FloatingIP: "8.8.8.8",
	}
	fixture3 := make(map[string]string)
	fixture3["ansible_host"] = "127.0.0.3"
	fixture3["floating_ip"] = "8.8.8.8"
	fixture3["swarm_labels"] = "google-dns"

	data := make([]testFixture, 0, 3)
	data = append(data, testFixture{
		Node:     server1,
		HostVars: fixture1,
	})

	data = append(data, testFixture{
		Node:     server2,
		HostVars: fixture2,
	})

	data = append(data, testFixture{
		Node:     server3,
		HostVars: fixture3,
	})

	return data
}
