package migrations

import (
	"github.com/koko1123/flow-go-1/ledger"
)

func NoOpMigration(p []ledger.Payload) ([]ledger.Payload, error) {
	return p, nil
}
