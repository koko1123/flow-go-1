package inmem

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

type Params struct {
	enc EncodableParams
}

func (p Params) ChainID() (flow.ChainID, error) {
	return p.enc.ChainID, nil
}

func (p Params) SporkID() (flow.Identifier, error) {
	return p.enc.SporkID, nil
}

func (p Params) ProtocolVersion() (uint, error) {
	return p.enc.ProtocolVersion, nil
}
