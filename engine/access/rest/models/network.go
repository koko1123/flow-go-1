package models

import (
	"github.com/koko1123/flow-go-1/access"
)

func (t *NetworkParameters) Build(params *access.NetworkParameters) {
	t.ChainId = params.ChainID.String()
}
