package flow_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/utils/unittest"
)

func TestAccountPublicKey_MarshalJSON(t *testing.T) {
	accountPrivateKey, err := unittest.AccountKeyDefaultFixture()
	assert.NoError(t, err)

	accountKeyA := accountPrivateKey.PublicKey(42)

	encAccountKey, err := json.Marshal(&accountKeyA)
	assert.NoError(t, err)

	var accountKeyB flow.AccountPublicKey

	err = json.Unmarshal(encAccountKey, &accountKeyB)
	assert.NoError(t, err)

	assert.Equal(t, accountKeyA, accountKeyB)
}
