package interfaces

import (
	"context"
)

type PeeringServer interface {
	String() string
	Close(ctx context.Context) error
	Stat() PeeringServerStat
}

type PeeringServerStat struct {
	PeeringConnections uint64

	HandshakeAttempts  uint64
	HandshakeSucceeded uint64

	PeerMessageReceived uint64
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
	PeeringAttempts  uint64
	PeeringConnected uint64

	HandshakeAttempts  uint64
	HandshakeSucceeded uint64

	PeerMessageReceived uint64
}

func (stat PeeringClientStat) String() string {
	return toJSON(stat)
}
