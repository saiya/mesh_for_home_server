package ingress

import (
	"io"
	"net/http"

	"github.com/saiya/mesh_for_home_server/config"
)

func httpHandler(
	config *config.HTTPIngressConfig,
	defaultHandler func(http.ResponseWriter, *http.Request) error,
) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		err := httpHandlerImpl(w, r, config, defaultHandler)
		if err != nil {
			w.WriteHeader(500)
			_, err = io.WriteString(w, "Internal Server Error")
			debugIfError("Failed to write 500 body", err)
		}
	})
}

func httpHandlerImpl(
	w http.ResponseWriter, r *http.Request,
	config *config.HTTPIngressConfig,
	defaultHandler func(http.ResponseWriter, *http.Request) error,
) error {
	handled, err := probeHandler(config, w, r)
	if handled {
		return err
	}

	return defaultHandler(w, r)
}

func probeHandler(config *config.HTTPIngressConfig, w http.ResponseWriter, r *http.Request) (bool, error) {
	if config.Probe == nil {
		return false, nil
	}

	if r.Method == "GET" && (config.Probe.Host == "" || r.Host == config.Probe.Host) && r.URL.Path == config.Probe.Path {
		w.WriteHeader(200)
		_, err := io.WriteString(w, "OK")
		return true, err
	}
	return false, nil
}
