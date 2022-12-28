package routetable_test

import (
	"context"
	"testing"
	"time"

	"github.com/saiya/mesh_for_home_server/config"
	"github.com/saiya/mesh_for_home_server/interfaces"
	"github.com/saiya/mesh_for_home_server/peering/proto/generated"
	"github.com/saiya/mesh_for_home_server/router/routetable"
	"github.com/stretchr/testify/assert"
)

func TestEmpty(t *testing.T) {
	ctx := context.Background()
	rt := routetable.NewHTTPRoutingTable()
	defer rt.Close(ctx)

	_, err := rt.Route(ctx, time.Now(), makeHTTPRequest("example.com"))
	assert.ErrorIs(t, err, interfaces.ErrNoRoute)
}

func TestMatch(t *testing.T) {
	ctx := context.Background()
	rt := routetable.NewHTTPRoutingTable()
	defer rt.Close(ctx)

	node1 := config.GenerateNodeID("test1")
	rt.Update(ctx, node1, time.Now().Add(3*time.Hour), makeHTTPAd("*.example.com"))

	node2 := config.GenerateNodeID("test2")
	rt.Update(ctx, node2, time.Now().Add(3*time.Hour), makeHTTPAd("sub.example.com"))

	dest, err := rt.Route(ctx, time.Now(), makeHTTPRequest("test.example.com"))
	assert.NoError(t, err)
	assert.Equal(t, node1, dest)

	dest, err = rt.Route(ctx, time.Now(), makeHTTPRequest("sub.example.com"))
	assert.NoError(t, err)
	assert.Equal(t, node2, dest)
}

func TestStale(t *testing.T) {
	ctx := context.Background()
	rt := routetable.NewHTTPRoutingTable()
	defer rt.Close(ctx)

	node1 := config.GenerateNodeID("test1")
	rt.Update(ctx, node1, time.Now().Add(-3*time.Hour), makeHTTPAd("example.com"))

	_, err := rt.Route(ctx, time.Now(), makeHTTPRequest("example.com"))
	assert.ErrorIs(t, err, interfaces.ErrNoRoute)
}

func TestRoundRobin(t *testing.T) {
	ctx := context.Background()
	rt := routetable.NewHTTPRoutingTable()
	defer rt.Close(ctx)

	nodes := make([]config.NodeID, 4)
	for i := range nodes {
		nodes[i] = config.GenerateNodeID("test")
		rt.Update(ctx, nodes[i], time.Now().Add(3*time.Hour), makeHTTPAd("example.com"))
	}

	staleNode := config.GenerateNodeID("test-stale")
	rt.Update(ctx, staleNode, time.Now().Add(-3*time.Hour), makeHTTPAd("example.com"))

	counter := make(map[config.NodeID]int)
	for i := 0; i < 128; i++ {
		dest, err := rt.Route(ctx, time.Now(), makeHTTPRequest("example.com"))
		assert.NoError(t, err)
		counter[dest]++
	}

	for _, node := range nodes {
		assert.NotZero(t, counter[node])
	}
	assert.Zero(t, counter[staleNode])
}

func makeHTTPAd(hostpatterns ...string) *generated.HttpAdvertisement {
	return &generated.HttpAdvertisement{HostnameMatchers: hostpatterns}
}

func makeHTTPRequest(hostname string) *generated.HttpRequestStart {
	return &generated.HttpRequestStart{
		Hostname: hostname,
	}
}
