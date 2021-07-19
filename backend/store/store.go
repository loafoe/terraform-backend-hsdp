package store

import (
	"errors"

	"github.com/philips-labs/terraform-backend-hsdp/backend/types"
)

// ErrNotFound item not found
var ErrNotFound = errors.New("resource not found")

// Store store interface
type Store interface {
	Init() error

	// state
	GetState(ref string, version ...string) (state map[string]interface{}, encrypted bool, err error)
	PutState(ref string, state, metadata map[string]interface{}, encrypted bool, version ...string) error
	DeleteState(ref string) error

	// lock
	GetLock(ref string) (lock *types.Lock, err error)
	PutLock(ref string, lock types.Lock) error
	DeleteLock(ref string) error

	// versioning
	List(ref string) ([]string, error)
	Restore(ref, version string) error
	Keep(last int) error
}
