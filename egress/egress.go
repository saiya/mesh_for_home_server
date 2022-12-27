package egress

import (
	"context"
	"fmt"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/egress/handler"
	"github.com/saiya/mesh_for_home_server/egress/handler/httphandler"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

func StartEgress(c *config.EgressConfigs, router interfaces.Router) (interfaces.AdvertisementProvider, []interfaces.MessageHandler, error) {
	httpHandler := httphandler.NewHttpHandler(router)
	handlers := []interfaces.MessageHandler{
		handler.NewPingHandler(router),
		httpHandler,
	}

	adExpreFunc := newADExpireFunc(c)
	advertiser := func(ctx context.Context) (interfaces.Advertisement, error) {
		return &generated.Advertisement{
			ExpireAt: adExpreFunc().Unix(),
			Http:     httpHandler.Advertise(),
		}, nil
	}
	if c == nil {
		return advertiser, handlers, nil
	}

	for i := range c.HTTP {
		if err := httpHandler.AddEgress(&c.HTTP[i]); err != nil {
			return advertiser, handlers, fmt.Errorf("failed to start HTTP(S) egress [%d]: %w", i, err)
		}
	}

	return advertiser, handlers, nil
}

func newADExpireFunc(c *config.EgressConfigs) func() time.Time {
	ttl := config.AdvertiseIntervalDefault
	if c != nil && c.AdvertiseIntervalSec > 0 {
		ttl = time.Duration(c.AdvertiseIntervalSec) * time.Second
	}
	ttl += config.AdvertiseTtlMargin

	return func() time.Time {
		return time.Now().Add(ttl)
	}
}
