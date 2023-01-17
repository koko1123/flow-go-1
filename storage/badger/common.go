package badger

import (
	"errors"
	"fmt"

	"github.com/dgraph-io/badger/v3"

	"github.com/koko1123/flow-go-1/storage"
)

func handleError(err error, t interface{}) error {
	if err != nil {
		if errors.Is(err, badger.ErrKeyNotFound) {
			return storage.ErrNotFound
		}

		return fmt.Errorf("could not retrieve %T: %w", t, err)
	}
	return nil
}
