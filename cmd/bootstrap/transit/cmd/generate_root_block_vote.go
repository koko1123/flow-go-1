package cmd

import (
	"encoding/json"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/koko1123/flow-go-1/cmd"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/model"
	hotstuffSig "github.com/koko1123/flow-go-1/consensus/hotstuff/signature"
	"github.com/koko1123/flow-go-1/consensus/hotstuff/verification"
	"github.com/koko1123/flow-go-1/model/bootstrap"
	"github.com/koko1123/flow-go-1/model/encodable"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/module/local"
	"github.com/koko1123/flow-go-1/utils/io"
)

var generateVoteCmd = &cobra.Command{
	Use:   "generate-root-block-vote",
	Short: "Generate root block vote",
	Run:   generateVote,
}

func init() {
	rootCmd.AddCommand(generateVoteCmd)
}

func generateVote(c *cobra.Command, args []string) {
	log.Info().Msg("generating root block vote")

	nodeIDString, err := readNodeID()
	if err != nil {
		log.Fatal().Err(err).Msg("could not read node ID")
	}

	nodeID, err := flow.HexStringToIdentifier(nodeIDString)
	if err != nil {
		log.Fatal().Err(err).Msg("could not parse node ID")
	}

	nodeInfo, err := cmd.LoadPrivateNodeInfo(flagBootDir, nodeID)
	if err != nil {
		log.Fatal().Err(err).Msg("could not load private node info")
	}

	// load DKG private key
	path := fmt.Sprintf(bootstrap.PathRandomBeaconPriv, nodeID)
	data, err := io.ReadFile(filepath.Join(flagBootDir, path))
	if err != nil {
		log.Fatal().Err(err).Msg("could not read DKG private key file")
	}

	var randomBeaconPrivKey encodable.RandomBeaconPrivKey
	err = json.Unmarshal(data, &randomBeaconPrivKey)
	if err != nil {
		log.Fatal().Err(err).Msg("could not unmarshal DKG private key data")
	}

	stakingPrivKey := nodeInfo.StakingPrivKey.PrivateKey
	identity := &flow.Identity{
		NodeID:        nodeID,
		Address:       nodeInfo.Address,
		Role:          nodeInfo.Role,
		Weight:        flow.DefaultInitialWeight,
		StakingPubKey: stakingPrivKey.PublicKey(),
		NetworkPubKey: nodeInfo.NetworkPrivKey.PrivateKey.PublicKey(),
	}

	me, err := local.New(identity, nodeInfo.StakingPrivKey.PrivateKey)
	if err != nil {
		log.Fatal().Err(err).Msg("creating local signer abstraction failed")
	}

	beaconKeyStore := hotstuffSig.NewStaticRandomBeaconSignerStore(randomBeaconPrivKey)
	signer := verification.NewCombinedSigner(me, beaconKeyStore)

	path = filepath.Join(flagBootDir, bootstrap.PathRootBlockData)
	data, err = io.ReadFile(path)
	if err != nil {
		log.Fatal().Err(err).Msg("could not read root block file")
	}

	var rootBlock flow.Block
	err = json.Unmarshal(data, &rootBlock)
	if err != nil {
		log.Fatal().Err(err).Msg("could not unmarshal root block data")
	}

	vote, err := signer.CreateVote(model.GenesisBlockFromFlow(rootBlock.Header))
	if err != nil {
		log.Fatal().Err(err).Msg("could not load private node info")
	}

	voteFile := fmt.Sprintf(bootstrap.PathNodeRootBlockVote, nodeID)

	if err = io.WriteJSON(filepath.Join(flagBootDir, voteFile), vote); err != nil {
		log.Fatal().Err(err).Msg("could not write vote to file")
	}

	log.Info().Msgf("node %v successfully generated vote file for block %v", nodeID, rootBlock.ID())
}
