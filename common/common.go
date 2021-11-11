package common

import (
	"context"
	"git.fd.io/govpp.git/api"
	"git.fd.io/govpp.git/core"
)

// Common features for cli commands
type Common struct {
	Root     string
	Connection *core.Connection
	Channel  api.Channel
	Context  context.Context
	Cancel   context.CancelFunc
}
