package events

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/messages"
)

type SyncedBlock struct {
	OriginID flow.Identifier
	Block    messages.UntrustedBlock
}

type SyncedClusterBlock struct {
	OriginID flow.Identifier
	Block    messages.UntrustedClusterBlock
}
