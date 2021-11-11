package common

import (
	"context"
	"git.fd.io/govpp.git/api"

)

// Common features for cli commands
type Common struct {
	Root     string
	channel  api.Channel
	Context  context.Context
	Cancel   context.CancelFunc
}
