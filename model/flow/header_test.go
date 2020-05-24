package flow_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vmihailenco/msgpack/v4"

	"github.com/dapperlabs/flow-go/crypto"
	"github.com/dapperlabs/flow-go/model/encoding/rlp"
	"github.com/dapperlabs/flow-go/model/flow"
	"github.com/dapperlabs/flow-go/utils/unittest"
)

func TestHeaderEncodingJSON(t *testing.T) {
	header := unittest.BlockHeaderFixture()
	headerID := header.ID()
	data, err := json.Marshal(header)
	require.NoError(t, err)
	var decoded flow.Header
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)
	decodedID := decoded.ID()
	assert.Equal(t, headerID, decodedID)
	assert.Equal(t, header, decoded)
}

func TestHeaderFingerprint(t *testing.T) {
	header := unittest.BlockHeaderFixture()
	headerID := header.ID()
	data := header.Fingerprint()
	var decoded struct {
		ChainID        string
		ParentID       flow.Identifier
		Height         uint64
		PayloadHash    flow.Identifier
		Timestamp      uint64
		View           uint64
		ParentVoterIDs []flow.Identifier
		ParentVoterSig crypto.Signature
		ProposerID     flow.Identifier
	}
	rlp.NewEncoder().MustDecode(data, &decoded)
	decHeader := flow.Header{
		ChainID:        decoded.ChainID,
		ParentID:       decoded.ParentID,
		Height:         decoded.Height,
		PayloadHash:    decoded.PayloadHash,
		Timestamp:      time.Unix(0, int64(decoded.Timestamp)).UTC(),
		View:           decoded.View,
		ParentVoterIDs: decoded.ParentVoterIDs,
		ParentVoterSig: decoded.ParentVoterSig,
		ProposerID:     decoded.ProposerID,
		ProposerSig:    header.ProposerSig, // since this field is not encoded/decoded, just set it to the original
		// value to pass test
	}
	decodedID := decHeader.ID()
	assert.Equal(t, headerID, decodedID)
	assert.Equal(t, header, decHeader)
}

func TestHeaderEncodingMsgpack(t *testing.T) {
	header := unittest.BlockHeaderFixture()
	headerID := header.ID()
	data, err := msgpack.Marshal(header)
	require.NoError(t, err)
	var decoded flow.Header
	err = msgpack.Unmarshal(data, &decoded)
	require.NoError(t, err)
	decodedID := decoded.ID()
	assert.Equal(t, headerID, decodedID)
	assert.Equal(t, header, decoded)
}

func TestNonUTCTimestampSameHashAsUTC(t *testing.T) {
	header := unittest.BlockHeaderFixture()
	headerID := header.ID()
	loc := time.FixedZone("UTC-8", -8*60*60)
	header.Timestamp = header.Timestamp.In(loc)
	checkedID := header.ID()
	assert.Equal(t, headerID, checkedID)
}
