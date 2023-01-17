package dkg

import (
	"github.com/koko1123/flow-go-1/crypto/hash"
	"github.com/koko1123/flow-go-1/module/signature"
)

// NewDKGMessageHasher returns a hasher for signing and verifying DKG broadcast
// messages.
func NewDKGMessageHasher() hash.Hasher {
	return signature.NewBLSHasher(signature.DKGMessageTag)
}
