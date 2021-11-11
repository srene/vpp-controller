package main

import (
	"context"
	"git.fd.io/govpp.git/api"
)

// Common features for cli commands
type Common struct {
	Routes map[string]string
	Stream   api.Stream
	Channel  api.Channel
	Context  context.Context
	Cancel   context.CancelFunc
}
