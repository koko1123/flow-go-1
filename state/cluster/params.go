package cluster

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

// Params contains constant information about this cluster state.
type Params interface {

	// ChainID returns the chain ID for this cluster.
	ChainID() (flow.ChainID, error)
}
