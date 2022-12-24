package ingress

import (
	"fmt"
	"net/http"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/ingress/forwarder"
	"github.com/saiya/mesh_for_home_server/ingress/httpingress"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

func StartIngress(c *config.IngressConfigs, router interfaces.Router, httpExecutor http.RoundTripper) ([]interfaces.Ingress, []interfaces.Forwarder, error) {
	ingresses := []interfaces.Ingress{}
	forwarders := []interfaces.Forwarder{forwarder.NewPingForwarder(router)}
	if c == nil {
		return ingresses, forwarders, nil
	}

	var httpForwarder interfaces.HTTPForwarder
	for i := range c.HTTP {
		if httpForwarder == nil {
			httpForwarder = forwarder.NewHTTPForwarder(router)
			forwarders = append(forwarders, httpForwarder)
		}

		cfg := &c.HTTP[i]
		httpIngress, err := httpingress.NewHTTPIngress(cfg, httpingress.NewDefaultHTTPHandler(httpForwarder))
		if err != nil {
			return ingresses, forwarders, fmt.Errorf("failed to start HTTP(S) ingress [%d]: %w", i, err)
		}
		ingresses = append(ingresses, httpIngress)
	}

	return ingresses, forwarders, nil
}
