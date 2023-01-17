package hotstuff

import (
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module"
)

// EventLoop performs buffer and processing of incoming proposals and QCs.
type EventLoop interface {
	module.HotStuff

	// SubmitTrustedQC accepts QC for processing. QC will be dispatched on worker thread.
	// CAUTION: QC is trusted (_not_ validated again), as it's built by ourselves.
	SubmitTrustedQC(qc *flow.QuorumCertificate)
}
