package ingress

import "context"

type Ingress interface {
	String() string
	Close(ctx context.Context) error
}
