package interfaces

import (
	"context"
	"errors"
	"sync/atomic"
)

var ErrPeeringDown = errors.New("Peering connection not available")

type PeeringServer interface {
	String() string
	Close(ctx context.Context) error
	Stat() PeeringServerStat
}

type PeeringServerStat struct {
	PeeringConnections  uint64
	HandshakeAttempts   uint64
	HandshakeSucceeded  uint64
	PeerMessageReceived uint64
}

func (stat *PeeringServerStat) Clone() PeeringServerStat {
	return PeeringServerStat{
		PeeringConnections:  atomic.LoadUint64(&stat.PeeringConnections),
		HandshakeAttempts:   atomic.LoadUint64(&stat.HandshakeAttempts),
		HandshakeSucceeded:  atomic.LoadUint64(&stat.HandshakeSucceeded),
		PeerMessageReceived: atomic.LoadUint64(&stat.PeerMessageReceived),
	}
}

func (stat PeeringServerStat) String() string {
	return toJSON(stat)
}

type PeeringClient interface {
	String() string
	Close(ctx context.Context) error
	Stat() PeeringClientStat
}

type PeeringClientStat struct {
	PeeringAttempts     uint64
	PeeringConnected    uint64
	HandshakeAttempts   uint64
	HandshakeSucceeded  uint64
	PeerMessageReceived uint64
}

func (stat *PeeringClientStat) Clone() PeeringClientStat {
	return PeeringClientStat{
		PeeringAttempts:     atomic.LoadUint64(&stat.PeeringAttempts),
		PeeringConnected:    atomic.LoadUint64(&stat.PeeringConnected),
		HandshakeAttempts:   atomic.LoadUint64(&stat.HandshakeAttempts),
		HandshakeSucceeded:  atomic.LoadUint64(&stat.HandshakeSucceeded),
		PeerMessageReceived: atomic.LoadUint64(&stat.PeerMessageReceived),
	}
}

func (stat PeeringClientStat) String() string {
	return toJSON(stat)
}
