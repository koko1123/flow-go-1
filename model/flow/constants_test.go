package flow_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/koko1123/flow-go-1/model/flow"
)

func TestDomainTags(t *testing.T) {
	assert.Len(t, flow.TransactionDomainTag, flow.DomainTagLength)
}
