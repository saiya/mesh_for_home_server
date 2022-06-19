package interfaces

import "context"

type Egress interface {
	String() string
	Close(ctx context.Context) error
}
