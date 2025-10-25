package response

import (
	log "log/slog"
)

func BuildEmptyRepository(err error) string {
	// print empty inventory
	if err != nil {
		// --debug prints something
		log.Debug(err.Error())
	}

	return "{\"_meta\": {\"hostvars\": {}}}"
}
