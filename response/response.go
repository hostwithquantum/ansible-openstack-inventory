package response

import (
	log "github.com/sirupsen/logrus"
)

func BuildEmptyRepository(err error) string {
	// print empty inventory
	if err != nil {
		// --debug prints something
		log.Debug(err)
	}

	return "{\"_meta\": {\"hostvars\": {}}}"
}
