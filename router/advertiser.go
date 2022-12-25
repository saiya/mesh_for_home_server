package router

import (
	"context"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
)

type advertiser struct {
	advFn     interfaces.Advertiser
	outbounds *outbounds

	ctx   context.Context
	close context.CancelFunc
	timer time.Timer
}

const advertiseInitialDelay = 500 * time.Millisecond
const advertiseRetryDelay = 3000 * time.Second

func NewAdvertiser(advFn interfaces.Advertiser, outbounds *outbounds) *advertiser {
	ctx, close := context.WithCancel(context.Background())
	av := &advertiser{
		advFn:     advFn,
		outbounds: outbounds,

		ctx:   ctx,
		close: close,
		timer: *time.NewTimer(advertiseInitialDelay),
	}
	go av.run()
	return av
}

func (av *advertiser) Close(ctx context.Context) error {
	av.close()
	av.timer.Stop()
	return nil
}

func (av *advertiser) GenerateAdvertisement(ctx context.Context) (interfaces.Advertisement, error) {
	return av.advFn(ctx)
}

func (av *advertiser) run() {
MainLoop:
	for {
		select {
		case <-av.ctx.Done():
			break MainLoop
		case <-av.timer.C:
			nextRun := av.advertise(av.ctx)
			av.timer = *time.NewTimer(nextRun.Sub(time.Now()))
		}
	}
}

func (av *advertiser) advertise(ctx context.Context) time.Time {
	packet, err := av.GenerateAdvertisement(ctx)
	if err != nil {
		logger.GetFrom(ctx).Errorw("Failed to generate advertisement: "+err.Error(), "err", err)
		return time.Now().Add(advertiseRetryDelay)
	}

	logger.GetFrom(ctx).Debugw("Broadcasting advertisement...")
	err = av.outbounds.Broadcast(ctx, &generated.PeerMessage{
		Message: &generated.PeerMessage_Advertisement{
			Advertisement: packet,
		},
	})
	if err != nil {
		logger.GetFrom(ctx).Infow("Failed to advertise to some peers (connection down?): "+err.Error(), "err", err)
	}
	return time.Unix(packet.ExpireAt, 0).Add(-config.AdvertiseTtlMargin)
}
