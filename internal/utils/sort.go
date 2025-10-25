package utils

import (
	"regexp"
	"sort"
	"strconv"

	"github.com/hostwithquantum/ansible-openstack-inventory/server"
)

// SortServersByName sorts servers by extracting the numeric portion from node-XXX.NAME format
func SortServersByName(servers []server.AnsibleServer) {
	// Regex to extract the number from node-XXX.NAME
	re := regexp.MustCompile(`node-(\d+)\.`)

	sort.Slice(servers, func(i, j int) bool {
		// Extract numbers from both server names
		numI := extractNumber(servers[i].Name, re)
		numJ := extractNumber(servers[j].Name, re)

		// Compare numerically
		return numI < numJ
	})
}

// extractNumber extracts the numeric portion from a server name
func extractNumber(name string, re *regexp.Regexp) int {
	matches := re.FindStringSubmatch(name)
	if len(matches) < 2 {
		return 0
	}

	num, err := strconv.Atoi(matches[1])
	if err != nil {
		return 0
	}

	return num
}
