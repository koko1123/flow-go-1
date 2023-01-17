package badger

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

type Params struct {
	state *State
}

func (p *Params) ChainID() (flow.ChainID, error) {
	return p.state.clusterID, nil
}
