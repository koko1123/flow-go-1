package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/koko1123/flow-go-1/cmd/bootstrap/run"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	"github.com/koko1123/flow-go-1/model/bootstrap"
	"github.com/koko1123/flow-go-1/model/dkg"
	"github.com/koko1123/flow-go-1/model/flow"
)

// constructRootQC constructs root QC based on root block, votes and dkg info
func constructRootQC(block *flow.Block, votes []*model.Vote, allNodes, internalNodes []bootstrap.NodeInfo, dkgData dkg.DKGData) *flow.QuorumCertificate {

	identities := bootstrap.ToIdentityList(allNodes)
	participantData, err := run.GenerateQCParticipantData(allNodes, internalNodes, dkgData)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate QC participant data")
	}

	qc, err := run.GenerateRootQC(block, votes, participantData, identities)
	if err != nil {
		log.Fatal().Err(err).Msg("generating root QC failed")
	}

	return qc
}

// NOTE: allNodes must be in the same order as when generating the DKG
func constructRootVotes(block *flow.Block, allNodes, internalNodes []bootstrap.NodeInfo, dkgData dkg.DKGData) {
	participantData, err := run.GenerateQCParticipantData(allNodes, internalNodes, dkgData)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to generate QC participant data")
	}

	votes, err := run.GenerateRootBlockVotes(block, participantData)
	if err != nil {
		log.Fatal().Err(err).Msg("generating votes for root block failed")
	}

	for _, vote := range votes {
		path := filepath.Join(bootstrap.DirnameRootBlockVotes, fmt.Sprintf(bootstrap.FilenameRootBlockVote, vote.SignerID))
		writeJSON(path, vote)
	}
}
