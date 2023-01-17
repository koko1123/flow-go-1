package utils

import (
	"github.com/koko1123/flow-go-1/crypto/hash"
	"github.com/koko1123/flow-go-1/module/signature"
)

// NewResultApprovalHasher generates and returns a hasher for signing
// and verification of result approvals
func NewResultApprovalHasher() hash.Hasher {
	h := signature.NewBLSHasher(signature.ResultApprovalTag)
	return h
}
