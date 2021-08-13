package server

const role_label string = "com.planetary-quantum.meta.role"

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
