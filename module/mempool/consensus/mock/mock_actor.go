package mock

import "github.com/koko1123/flow-go-1/model/flow"

// ExecForkActor allows to create a mock for the ExecForkActor callback
type ExecForkActor interface {
	OnExecFork([]*flow.IncorporatedResultSeal)
}
