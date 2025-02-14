package protocol

import (
	"errors"
	"fmt"

	"github.com/koko1123/flow-go-1/model/flow"
	"github.com/koko1123/flow-go-1/state"
)

var (
	// ErrNoPreviousEpoch is a sentinel error returned when a previous epoch is
	// queried from a snapshot within the first epoch after the root block.
	ErrNoPreviousEpoch = fmt.Errorf("no previous epoch exists")

	// ErrNextEpochNotSetup is a sentinel error returned when the next epoch
	// has not been set up yet.
	ErrNextEpochNotSetup = fmt.Errorf("next epoch has not yet been set up")

	// ErrEpochNotCommitted is a sentinel error returned when the epoch has
	// not been committed and information is queried that is only accessible
	// in the EpochCommitted phase.
	ErrEpochNotCommitted = fmt.Errorf("queried info from EpochCommit event before it was emitted")

	// ErrSealingSegmentBelowRootBlock is a sentinel error returned for queries
	// for a sealing segment below the root block.
	ErrSealingSegmentBelowRootBlock = fmt.Errorf("cannot query sealing segment below root block")

	// ErrClusterNotFound is a sentinel error returns for queries for a cluster
	ErrClusterNotFound = fmt.Errorf("could not find cluster")
)

type IdentityNotFoundError struct {
	NodeID flow.Identifier
}

func (e IdentityNotFoundError) Error() string {
	return fmt.Sprintf("identity not found (%x)", e.NodeID)
}

func IsIdentityNotFound(err error) bool {
	var errIdentityNotFound IdentityNotFoundError
	return errors.As(err, &errIdentityNotFound)
}

type InvalidBlockTimestampError struct {
	error
}

func (e InvalidBlockTimestampError) Unwrap() error {
	return e.error
}

func (e InvalidBlockTimestampError) Error() string {
	return e.error.Error()
}

func IsInvalidBlockTimestampError(err error) bool {
	var errInvalidTimestampError InvalidBlockTimestampError
	return errors.As(err, &errInvalidTimestampError)
}

func NewInvalidBlockTimestamp(msg string, args ...interface{}) error {
	return InvalidBlockTimestampError{
		error: fmt.Errorf(msg, args...),
	}
}

// InvalidServiceEventError indicates an invalid service event was processed.
type InvalidServiceEventError struct {
	error
}

func (e InvalidServiceEventError) Unwrap() error {
	return e.error
}

func IsInvalidServiceEventError(err error) bool {
	var errInvalidServiceEventError InvalidServiceEventError
	return errors.As(err, &errInvalidServiceEventError)
}

// NewInvalidServiceEventError returns an invalid service event error. Since all invalid
// service events indicate an invalid extension, the service event error is wrapped in
// the invalid extension error at construction.
func NewInvalidServiceEventError(msg string, args ...interface{}) error {
	return state.NewInvalidExtensionErrorf(
		"cannot extend state with invalid service event: %w",
		InvalidServiceEventError{
			error: fmt.Errorf(msg, args...),
		},
	)
}
