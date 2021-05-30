package server_test

import (
	"testing"

	"github.com/hostwithquantum/ansible-openstack-inventory/server"
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

		if server.IsWorker(node) {
			t.Error("Shouldn't be a worker")
		}
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

		if !server.IsWorker(node) {
			t.Error("Should be a worker")
		}
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

		if !server.IsManager(node) {
			t.Error("Should be a manager")
		}
	}
}
