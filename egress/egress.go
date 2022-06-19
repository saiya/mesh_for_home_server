package egress

import (
	"fmt"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress/httpegress"
	"github.com/saiya/mesh_for_home_server/interfaces"
)

func StartEgress(c *config.EgressConfig) ([]interfaces.Egress, error) {
	result := []interfaces.Egress{}
	for i := range c.HTTP {
		cfg := &c.HTTP[i]
		httpEgress, err := httpegress.NewHTTPEgress(cfg)
		if err != nil {
			return result, fmt.Errorf("failed to start HTTP(S) egress [%d]: %w", i, err)
		}
		result = append(result, httpEgress)
	}
	return result, nil
}
