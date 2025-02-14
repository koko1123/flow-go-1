package committees

import (
	"fmt"

	"github.com/koko1123/flow-go-1/consensus/hotstuff"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/model/flow/order"
	"github.com/koko1123/flow-go-1/state/protocol"
	"github.com/onflow/flow-go/crypto"
)

// NewStaticCommittee returns a new committee with a static participant set.
func NewStaticCommittee(participants flow.IdentityList, myID flow.Identifier, dkgParticipants map[flow.Identifier]flow.DKGParticipant, dkgGroupKey crypto.PublicKey) (*Static, error) {

	return NewStaticCommitteeWithDKG(participants, myID, staticDKG{
		dkgParticipants: dkgParticipants,
		dkgGroupKey:     dkgGroupKey,
	})
}

// NewStaticCommitteeWithDKG returns a new committee with a static participant set.
func NewStaticCommitteeWithDKG(participants flow.IdentityList, myID flow.Identifier, dkg protocol.DKG) (*Static, error) {
	valid := order.IdentityListCanonical(participants)
	if !valid {
		return nil, fmt.Errorf("participants %v is not in Canonical order", participants)
	}

	static := &Static{
		participants: participants,
		myID:         myID,
		dkg:          dkg,
	}
	return static, nil
}

// Static represents a committee with a static participant set. It is used for
// bootstrapping purposes.
type Static struct {
	participants flow.IdentityList
	myID         flow.Identifier
	dkg          protocol.DKG
}

func (s Static) Identities(_ flow.Identifier) (flow.IdentityList, error) {
	return s.participants, nil
}

func (s Static) Identity(_ flow.Identifier, participantID flow.Identifier) (*flow.Identity, error) {
	identity, ok := s.participants.ByNodeID(participantID)
	if !ok {
		return nil, fmt.Errorf("unknown partipant")
	}
	return identity, nil
}

func (s Static) LeaderForView(_ uint64) (flow.Identifier, error) {
	return flow.ZeroID, fmt.Errorf("invalid for static committee")
}

func (s Static) Self() flow.Identifier {
	return s.myID
}

func (s Static) DKG(_ flow.Identifier) (hotstuff.DKG, error) {
	return s.dkg, nil
}

type staticDKG struct {
	dkgParticipants map[flow.Identifier]flow.DKGParticipant
	dkgGroupKey     crypto.PublicKey
}

func (s staticDKG) Size() uint {
	return uint(len(s.dkgParticipants))
}

func (s staticDKG) GroupKey() crypto.PublicKey {
	return s.dkgGroupKey
}

// Index returns the index for the given node. Error Returns:
// protocol.IdentityNotFoundError if nodeID is not a valid DKG participant.
func (s staticDKG) Index(nodeID flow.Identifier) (uint, error) {
	participant, ok := s.dkgParticipants[nodeID]
	if !ok {
		return 0, protocol.IdentityNotFoundError{NodeID: nodeID}
	}
	return participant.Index, nil
}

// KeyShare returns the public key share for the given node. Error Returns:
// protocol.IdentityNotFoundError if nodeID is not a valid DKG participant.
func (s staticDKG) KeyShare(nodeID flow.Identifier) (crypto.PublicKey, error) {
	participant, ok := s.dkgParticipants[nodeID]
	if !ok {
		return nil, protocol.IdentityNotFoundError{NodeID: nodeID}
	}
	return participant.KeyShare, nil
}
