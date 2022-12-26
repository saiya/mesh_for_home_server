package router

import (
	"context"
	"sync"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/logger"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
	"github.com/saiya/mesh_for_home_server/router/routetable"
)

type router struct {
	nodeID config.NodeID

	rt         routetable.RouteTable
	advertiser *advertiser
	advFnLock  sync.Mutex
	advFn      interfaces.AdvertisementProvider

	outbounds *outbounds
	inbound   *inbound
}

func NewRouter(hostname string) interfaces.Router {
	outbounds := newOutbounds()
	router := &router{
		nodeID: config.GenerateNodeID(hostname),

		rt: routetable.NewRouteTable(),

		outbounds: outbounds,
		inbound:   newInbound(),
	}
	router.advertiser = newAdvertiser(router.GenerateAdvertisement, outbounds)

	router.inbound.Register(func(ctx context.Context, from config.NodeID, msg interfaces.Message) error {
		ad := msg.GetAdvertisement()
		if ad != nil {
			router.rt.Update(ctx, from, ad)
		}
		return nil
	})

	logger.Get().Infow("Router initialized", "nodeID", router.nodeID)
	return router
}

func (r *router) Close(ctx context.Context) error {
	return r.rt.Close(ctx)
}

func (r *router) NodeID() config.NodeID {
	return r.nodeID
}

func (r *router) SetAdvertisementProvider(advFn interfaces.AdvertisementProvider) {
	r.advFnLock.Lock()
	defer r.advFnLock.Unlock()

	r.advFn = advFn
}

func (r *router) GenerateAdvertisement(ctx context.Context) (interfaces.Advertisement, error) {
	var advFn interfaces.AdvertisementProvider

	func() {
		r.advFnLock.Lock()
		defer r.advFnLock.Unlock()
		advFn = r.advFn
	}()

	if advFn == nil {
		logger.GetFrom(ctx).Errorw("AdvertisementProvider not found, stubbing with empty AD...")
		return &generated.Advertisement{}, nil
	}
	return advFn(ctx)
}

func (r *router) Update(ctx context.Context, node config.NodeID, ad interfaces.Advertisement) {
	r.rt.Update(ctx, node, ad)
}

func (r *router) Route(ctx context.Context, request interfaces.Message) (config.NodeID, error) {
	return r.rt.Route(ctx, request)
}

func (r *router) RegisterSink(dest config.NodeID, callback interfaces.RouterShink) interfaces.RouterUnregister {
	if r.nodeID == dest {
		panic("Sink's destination nodeID must be different from this node's ID")
	}

	return r.outbounds.Register(dest, callback)
}

func (r *router) Listen(callback interfaces.RouterListener) interfaces.RouterUnregister {
	return r.inbound.Register(callback)
}

func (r *router) Deliver(parentCtx context.Context, from config.NodeID, dest config.NodeID, msg interfaces.Message) {
	ctx := logger.Wrap(parentCtx, "from", from, "dest", dest, "peer-msg", interfaces.MsgLogString(msg))
	if r.nodeID == dest {
		r.inbound.Invoke(ctx, from, msg)
	} else {
		err := r.outbounds.Deliver(ctx, dest, msg)
		if err != nil {
			logger.GetFrom(ctx).Warnw("failed to deliver message", "err", err)
		}
	}
}
