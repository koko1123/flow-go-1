package mocksiface_test

import (
	"github.com/koko1123/flow-go-1-sdk/access"
)

// This is a proxy for the real access.Client for mockery to use.
type Client interface {
	access.Client
}
