package state

import (
	"github.com/koko1123/flow-go-1/model/flow"
)

type View interface {
	NewChild() View
	MergeView(child View) error
	DropDelta() // drops all the delta changes
	RegisterUpdates() ([]flow.RegisterID, []flow.RegisterValue)
	AllRegisters() []flow.RegisterID
	Ledger
}

// Ledger is the storage interface used by the virtual machine to read and write register values.
//
// TODO Rename this to Storage
// and remove reference to flow.RegisterValue and use byte[]
type Ledger interface {
	Set(owner, key string, value flow.RegisterValue) error
	Get(owner, key string) (flow.RegisterValue, error)
	Touch(owner, key string) error
	Delete(owner, key string) error
}
