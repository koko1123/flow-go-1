package fvm

import (
	"github.com/koko1123/flow-go-1/fvm/environment"
	"github.com/koko1123/flow-go-1/storage"
)

// TODO(patrick): rm after https://github.com/onflow/flow-emulator/pull/229
// is merged and integrated.
type Blocks = environment.Blocks

// TODO(patrick): rm after https://github.com/onflow/flow-emulator/pull/229
// is merged and integrated.
func NewBlockFinder(storage storage.Headers) Blocks {
	return environment.NewBlockFinder(storage)
}
