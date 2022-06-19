package ingress

import (
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/ingress/httpingress"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

func StartIngress(c *config.IngressConfig, router interfaces.Router) ([]interfaces.Ingress, error) {
	result := []interfaces.Ingress{}
	for i := range c.HTTP {
		cfg := &c.HTTP[i]
		httpIngress, err := httpingress.NewHTTPIngress(cfg, httpingress.NewDefaultHTTPHandler(router.HTTP()))
		if err != nil {
			return result, fmt.Errorf("failed to start HTTP(S) ingress [%d]: %w", i, err)
		}
		result = append(result, httpIngress)
	}
	return result, nil
}
