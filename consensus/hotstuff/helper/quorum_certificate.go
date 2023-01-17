package helper

import (
	"math/rand"

	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func MakeQC(options ...func(*flow.QuorumCertificate)) *flow.QuorumCertificate {
	qc := flow.QuorumCertificate{
		View:          rand.Uint64(),
		BlockID:       unittest.IdentifierFixture(),
		SignerIndices: unittest.SignerIndicesFixture(3),
		SigData:       unittest.SignatureFixture(),
	}
	for _, option := range options {
		option(&qc)
	}
	return &qc
}

func WithQCBlock(block *model.Block) func(*flow.QuorumCertificate) {
	return func(qc *flow.QuorumCertificate) {
		qc.View = block.View
		qc.BlockID = block.BlockID
	}
}

func WithQCSigners(signerIndices []byte) func(*flow.QuorumCertificate) {
	return func(qc *flow.QuorumCertificate) {
		qc.SignerIndices = signerIndices
	}
}

func WithQCView(view uint64) func(*flow.QuorumCertificate) {
	return func(qc *flow.QuorumCertificate) {
		qc.View = view
	}
}
