package request

import (
	"fmt"

	"github.com/koko1123/flow-go-1/engine/access/rest/models"
	"github.com/koko1123/flow-go-1/engine/access/rest/util"
	"github.com/koko1123/flow-go-1/model/flow"
)

type ProposalKey flow.ProposalKey

func (p *ProposalKey) Parse(raw models.ProposalKey) error {
	var address Address
	err := address.Parse(raw.Address)
	if err != nil {
		return err
	}

	keyIndex, err := util.ToUint64(raw.KeyIndex)
	if err != nil {
		return fmt.Errorf("invalid key index: %w", err)
	}

	seqNumber, err := util.ToUint64(raw.SequenceNumber)
	if err != nil {
		return fmt.Errorf("invalid sequence number: %w", err)
	}

	*p = ProposalKey(flow.ProposalKey{
		Address:        address.Flow(),
		KeyIndex:       keyIndex,
		SequenceNumber: seqNumber,
	})
	return nil
}

func (p ProposalKey) Flow() flow.ProposalKey {
	return flow.ProposalKey(p)
}
