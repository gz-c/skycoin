package api

import (
	"net/http"

	wh "github.com/skycoin/skycoin/src/util/http"
)

// versionHandler returns the application version info
// URI: /api/v1/version
// Method: GET
func versionHandler(gateway Gatewayer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			wh.Error405(w)
			return
		}

		wh.SendJSONOr500(logger, w, gateway.GetBuildInfo())
	}
}
