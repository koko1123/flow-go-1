package utils

import (
	"github.com/koko1123/flow-go-1/module/signature"
	"github.com/onflow/flow-go/crypto/hash"
)

// NewResultApprovalHasher generates and returns a hasher for signing
// and verification of result approvals
func NewResultApprovalHasher() hash.Hasher {
	h := signature.NewBLSHasher(signature.ResultApprovalTag)
	return h
}
