package server_test

import (
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/server"
	"github.com/stretchr/testify/assert"
)

func Test_NoWorker(t *testing.T) {
	metaDataExamples := []map[string]string{
		{},
		{"com.planetary-quantum.meta.role": "manager"},
		{"com.planetary-quantum.meta.role": "unicorn"},
	}

	for _, md := range metaDataExamples {
		node := server.AnsibleServer{
			MetaData: md,
		}

		assert.False(t, server.IsWorker(node), "Shouldn't be a worker")
	}
}

func Test_Worker(t *testing.T) {
	metaDataExamples := []map[string]string{
		{"com.planetary-quantum.meta.role": "worker"},
	}

	for _, md := range metaDataExamples {
		node := server.AnsibleServer{
			MetaData: md,
		}
		assert.True(t, server.IsWorker(node), "Should be a worker")
	}
}

func Test_Manager(t *testing.T) {
	metaDataExamples := []map[string]string{
		{},
		{"com.planetary-quantum.meta.role": "manager"},
	}

	for _, md := range metaDataExamples {
		node := server.AnsibleServer{
			MetaData: md,
		}

		assert.True(t, server.IsManager(node), "Should be a manager")
	}
}

func Test_Group(t *testing.T) {
	metaDataExamples := []map[string]string{
		{"com.planetary-quantum.meta.customer_group": "blah"},
		{"com.planetary-quantum.meta.customer_group": "example-customer-group"},
	}

	for _, md := range metaDataExamples {
		node := server.AnsibleServer{
			MetaData: md,
		}

		_, err := server.GetGroup(node)
		assert.NoError(t, err)
	}
}
