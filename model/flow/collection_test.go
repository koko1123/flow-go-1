package flow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/koko1123/flow-go-1/model/encoding/rlp"
	"github.com/koko1123/flow-go-1/model/fingerprint"
	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestLightCollectionFingerprint(t *testing.T) {
	col := unittest.CollectionFixture(2)
	colID := col.ID()
	data := fingerprint.Fingerprint(col.Light())
	var decoded flow.LightCollection
	rlp.NewMarshaler().MustUnmarshal(data, &decoded)
	decodedID := decoded.ID()
	assert.Equal(t, colID, decodedID)
	assert.Equal(t, col.Light(), decoded)
}
