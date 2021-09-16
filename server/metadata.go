package server

import "fmt"

const role_label string = "com.planetary-quantum.meta.role"
const group_label string = "com.planetary-quantum.meta.customer_group"

func GetGroup(server AnsibleServer) (string, error) {
	if len(server.MetaData) == 0 {
		return "", fmt.Errorf("Server doesn't have metadata attached: %s", server.Name)
	}

	group, ok := server.MetaData[group_label]
	if !ok {
		return "", fmt.Errorf("Server doesn't have role attached: %s", server.Name)
	}

	return group, nil
}

// IsManager ...
func IsManager(server AnsibleServer) bool {
	// fallback: setups without metadata
	if len(server.MetaData) == 0 {
		return true
	}

	return isRole(server.MetaData, "manager")
}

// IsWorker ...
func IsWorker(server AnsibleServer) bool {
	return isRole(server.MetaData, "worker")
}

func isRole(md map[string]string, role string) bool {
	if len(md) == 0 {
		return false
	}

	node_role, ok := md[role_label]
	if !ok {
		return false
	}

	if node_role == role {
		return true
	}

	return false
}
