package interfaces

import "context"

type MessageHandler interface {
	Close(ctx context.Context) error
}
