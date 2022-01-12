package request

import (
	"fmt"
	"github.com/onflow/flow-go/engine/access/rest/models"
	"github.com/onflow/flow-go/model/flow"
)

const signatureLength = 128

type TransactionSignature flow.TransactionSignature

func (s *TransactionSignature) Parse(sig models.TransactionSignature) error {
	var address Address
	err := address.Parse(sig.Address)
	if err != nil {
		return err
	}

	sigIndex, err := toUint64(sig.SignerIndex)
	if err != nil {
		return fmt.Errorf("invalid signer index: %v", err)
	}

	keyIndex, err := toUint64(sig.KeyIndex)
	if err != nil {
		return fmt.Errorf("invalid key index: %v", err)
	}

	var signature Signature
	err = signature.Parse(sig.Signature)
	if err != nil {
		return fmt.Errorf("invalid signature: %v", err)
	}

	*s = TransactionSignature(flow.TransactionSignature{
		Address:     address.Flow(),
		SignerIndex: int(sigIndex),
		KeyIndex:    keyIndex,
		Signature:   signature,
	})

	return nil
}

func (s TransactionSignature) Flow() flow.TransactionSignature {
	return flow.TransactionSignature(s)
}

type TransactionSignatures []TransactionSignature

func (t *TransactionSignatures) Parse(rawSigs []models.TransactionSignature) error {
	signatures := make([]TransactionSignature, len(rawSigs))
	for i, sig := range rawSigs {
		var signature TransactionSignature
		err := signature.Parse(sig)
		if err != nil {
			return err
		}
		signatures[i] = signature
	}
	return nil
}

func (t TransactionSignatures) Flow() []flow.TransactionSignature {
	sigs := make([]flow.TransactionSignature, len(t))
	for i, sig := range t {
		sigs[i] = sig.Flow()
	}
	return sigs
}

type Signature []byte

func (s *Signature) Parse(raw string) error {
	signatureBytes, err := fromBase64(raw)
	if err != nil {
		return fmt.Errorf("invalid encoding")
	}

	*s = signatureBytes
	return nil
}

func (s Signature) Flow() []byte {
	return s
}
