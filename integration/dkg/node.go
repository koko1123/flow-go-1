package dkg

import (
	"crypto"
	"testing"

	sdk "github.com/koko1123/flow-go-1-sdk"
	sdkcrypto "github.com/koko1123/flow-go-1-sdk/crypto"
	"github.com/koko1123/flow-go-1/engine/consensus/dkg"
	testmock "github.com/koko1123/flow-go-1/engine/testutil/mock"
	"github.com/koko1123/flow-go-1/model/bootstrap"
	"github.com/koko1123/flow-go-1/model/flow"
	protocolmock "github.com/koko1123/flow-go-1/state/protocol/mock"
	"github.com/koko1123/flow-go-1/storage"
	"github.com/koko1123/flow-go-1/utils/unittest/mocks"
)

type nodeAccount struct {
	netID          *flow.Identity
	privKey        crypto.PrivateKey
	accountKey     *sdk.AccountKey
	accountID      string
	accountAddress sdk.Address
	accountSigner  sdkcrypto.Signer
	accountInfo    *bootstrap.NodeMachineAccountInfo
}

// node is an in-process node that only contains the engines relevant to DKG,
// ie. MessagingEngine and ReactorEngine
type node struct {
	testmock.GenericNode
	account           *nodeAccount
	dkgContractClient *DKGClientWrapper
	dkgState          storage.DKGState
	safeBeaconKeys    storage.SafeBeaconKeys
	messagingEngine   *dkg.MessagingEngine
	reactorEngine     *dkg.ReactorEngine
}

func (n *node) Ready() {
	<-n.messagingEngine.Ready()
	<-n.reactorEngine.Ready()
}

func (n *node) Done() {
	<-n.messagingEngine.Done()
	<-n.reactorEngine.Done()
	// close database otherwise hitting "too many file open"
	_ = n.PublicDB.Close()
	_ = n.SecretsDB.Close()
}

// setEpochs configures the mock state snapthost at firstBlock to return the
// desired current and next epochs
func (n *node) setEpochs(t *testing.T, currentSetup flow.EpochSetup, nextSetup flow.EpochSetup, firstBlock *flow.Header) {

	currentEpoch := new(protocolmock.Epoch)
	currentEpoch.On("Counter").Return(currentSetup.Counter, nil)
	currentEpoch.On("InitialIdentities").Return(currentSetup.Participants, nil)
	currentEpoch.On("DKGPhase1FinalView").Return(currentSetup.DKGPhase1FinalView, nil)
	currentEpoch.On("DKGPhase2FinalView").Return(currentSetup.DKGPhase2FinalView, nil)
	currentEpoch.On("DKGPhase3FinalView").Return(currentSetup.DKGPhase3FinalView, nil)
	currentEpoch.On("FinalView").Return(currentSetup.FinalView, nil)
	currentEpoch.On("FirstView").Return(currentSetup.FirstView, nil)
	currentEpoch.On("RandomSource").Return(nextSetup.RandomSource, nil)

	nextEpoch := new(protocolmock.Epoch)
	nextEpoch.On("Counter").Return(nextSetup.Counter, nil)
	nextEpoch.On("InitialIdentities").Return(nextSetup.Participants, nil)
	nextEpoch.On("RandomSource").Return(nextSetup.RandomSource, nil)
	nextEpoch.On("DKG").Return(nil, nil) // no error means didn't run into EECC
	nextEpoch.On("FirstView").Return(nextSetup.FirstView, nil)
	nextEpoch.On("FinalView").Return(nextSetup.FinalView, nil)

	epochQuery := mocks.NewEpochQuery(t, currentSetup.Counter)
	epochQuery.Add(currentEpoch)
	epochQuery.Add(nextEpoch)
	snapshot := new(protocolmock.Snapshot)
	snapshot.On("Epochs").Return(epochQuery)
	snapshot.On("Phase").Return(flow.EpochPhaseStaking, nil)
	snapshot.On("Head").Return(firstBlock, nil)
	state := new(protocolmock.MutableState)
	state.On("AtBlockID", firstBlock.ID()).Return(snapshot)
	state.On("Final").Return(snapshot)
	n.GenericNode.State = state
	n.reactorEngine.State = state
}
