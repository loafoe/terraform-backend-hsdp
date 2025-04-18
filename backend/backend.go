package backend

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	gocrypto "github.com/bhoriuchi/go-crypto"
	"github.com/loafoe/terraform-backend-hsdp/backend/store"
	"github.com/loafoe/terraform-backend-hsdp/backend/types"
)

// Options backend options
type Options struct {
	EncryptionKey   interface{}
	Logger          func(level, message string, err error)
	GetRefFunc      interface{}
	GetEncryptFunc  interface{}
	GetMetadataFunc func(state map[string]interface{}) map[string]interface{}
}

// NewBackend creates a new backend
func NewBackend(store store.Store, opts ...*Options) *Backend {
	backend := Backend{
		initialized: false,
		store:       store,
	}

	if len(opts) > 0 {
		backend.options = opts[0]
	}
	if backend.options == nil {
		backend.options = &Options{}
	}
	if backend.options.Logger == nil {
		backend.options.Logger = func(level, message string, err error) {}
	}
	if backend.options.GetMetadataFunc == nil {
		backend.options.GetMetadataFunc = func(state map[string]interface{}) map[string]interface{} {
			return map[string]interface{}{}
		}
	}

	return &backend
}

// Backend a terraform http backend
type Backend struct {
	initialized bool
	store       store.Store
	options     *Options
}

// Init initializes the backend
func (c *Backend) Init() error {
	if !c.initialized {
		return c.store.Init()
	}
	return nil
}

// gets the encryption key
func (c *Backend) getEncryptionKey() []byte {
	switch key := c.options.EncryptionKey; key.(type) {
	case []byte:
		return key.([]byte)
	case func() []byte:
		return key.(func() []byte)()
	}
	return nil
}

// gets the state ref
func (c *Backend) getRef(r *http.Request) (string, error) {
	switch refFunc := c.options.GetRefFunc; refFunc.(type) {
	case func(r *http.Request) (string, error):
		return refFunc.(func(r *http.Request) (string, error))(r)
	}
	return r.URL.Query().Get("ref"), nil
}

// gets the encrypt state setting
func (c *Backend) getEncrypt(r *http.Request) bool {
	switch encFunc := c.options.GetRefFunc; encFunc.(type) {
	case func(r *http.Request) bool:
		return encFunc.(func(r *http.Request) bool)(r)
	}
	// Encrypt by default
	return true
}

// decrypts the encrypted state
func (c *Backend) decryptState(encryptedState interface{}) (map[string]interface{}, error) {
	key := c.getEncryptionKey()
	if len(key) == 0 {
		return nil, fmt.Errorf("failed to get backend encryption key")
	}

	s := types.EncryptedState{}
	if err := toInterface(encryptedState, &s); err != nil {
		return nil, err
	}

	data, err := base64.StdEncoding.DecodeString(s.EncryptedData)
	if err != nil {
		return nil, err
	}

	decryptedData, err := gocrypto.Decrypt(key, data)
	if err != nil {
		return nil, err
	}

	var state map[string]interface{}
	if err := json.Unmarshal(decryptedData, &state); err != nil {
		return nil, err
	}

	return state, nil
}

// encrypts the state
func (c *Backend) encryptState(state interface{}) (map[string]interface{}, error) {
	key := c.getEncryptionKey()
	if len(key) == 0 {
		return nil, fmt.Errorf("failed to get backend encryption key")
	}

	j, err := json.Marshal(state)
	if err != nil {
		return nil, err
	}

	encryptedData, err := gocrypto.Encrypt(key, j)
	if err != nil {
		return nil, err
	}

	var encryptedState map[string]interface{}
	s := types.EncryptedState{
		EncryptedData: base64.StdEncoding.EncodeToString(encryptedData),
	}

	if err := toInterface(s, &encryptedState); err != nil {
		return nil, err
	}

	return encryptedState, nil
}

// determines if the state can be locked
func (c *Backend) canLock(w http.ResponseWriter, _ *http.Request, ref, id string) bool {
	lock, err := c.store.GetLock(ref)
	if err != nil {
		if err == store.ErrNotFound {
			return true
		}

		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get lock from state store for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return false
	}

	if lock.ID == id {
		return true
	}

	c.options.Logger(
		"debug",
		fmt.Sprintf("terraform state locked by another process for ref: %s", ref),
		nil,
	)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusLocked)
	_ = json.NewEncoder(w).Encode(lock)
	return false
}

// HandleGetState gets the state requested
func (c *Backend) HandleGetState(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.options.Logger(
		"debug",
		fmt.Sprintf("getting terraform state for ref: %s", ref),
		nil,
	)
	// get the state
	state, encrypted, err := c.store.GetState(ref)
	if err != nil {
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get terraform state for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// decrypt
	if encrypted {
		decryptedState, err := c.decryptState(state)
		if err != nil {
			c.options.Logger(
				"error",
				fmt.Sprintf("failed decrypt terraform state for ref: %s", ref),
				err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		state = decryptedState
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(state)
}

// HandleLockState locks the state
func (c *Backend) HandleLockState(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.options.Logger(
		"debug",
		fmt.Sprintf("locking terraform state for ref %s", ref),
		nil,
	)

	// decode body
	var lock types.Lock
	if err := json.NewDecoder(r.Body).Decode(&lock); err != nil {
		c.options.Logger(
			"debug",
			fmt.Sprintf("error decoding LOCK request body for ref %s", ref),
			nil,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if state can be locked
	if !c.canLock(w, r, ref, lock.ID) {
		return
	}

	// attempt to put the lock
	if err := c.store.PutLock(ref, lock); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to set lock for ref %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleUnlockState unlocks the state
func (c *Backend) HandleUnlockState(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.options.Logger(
		"debug",
		fmt.Sprintf("unlocking terraform state for ref %s", ref),
		nil,
	)

	// decode body
	var lock types.Lock
	if err := json.NewDecoder(r.Body).Decode(&lock); err != nil && err != io.EOF {
		c.options.Logger(
			"error",
			fmt.Sprintf("error decoding UNLOCK request body for ref %s", ref),
			err,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// check if state can be locked
	if !c.canLock(w, r, ref, lock.ID) {
		return
	}

	// attempt to delete the lock
	if err := c.store.DeleteLock(ref); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to delete lock for ref %s", ref),
			err,
		)

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (c *Backend) getVersion(now time.Time) string {
	return now.Format("20060102150405")
}

// HandleUpdateState updates the state
func (c *Backend) HandleUpdateState(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	encrypt := c.getEncrypt(r)
	id := r.URL.Query().Get("ID")

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.options.Logger(
		"debug",
		fmt.Sprintf("setting terraform state for ref %s", ref),
		nil,
	)

	// decode body
	var state map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&state); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("error decoding request body for ref %s", ref),
			err,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !c.canLock(w, r, ref, id) {
		return
	}

	// get metadata using a metadata processor
	metadata := c.options.GetMetadataFunc(state)

	// encrypt if specified
	if encrypt {
		encryptedState, err := c.encryptState(state)
		if err != nil {
			c.options.Logger(
				"error",
				fmt.Sprintf("failed encrypt terraform state for ref: %s", ref),
				err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		state = encryptedState
	}

	// set the state on the backend
	if err := c.store.PutState(ref, state, metadata, encrypt); err != nil {
		c.options.Logger(
			"debug",
			fmt.Sprintf("error updating terraform state for ref %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// write a version
	_ = c.store.PutState(ref, state, metadata, encrypt, c.getVersion(time.Now()))

	w.WriteHeader(http.StatusOK)
}

// HandleDeleteState deletes the state
func (c *Backend) HandleDeleteState(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	id := r.URL.Query().Get("ID")

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	c.options.Logger(
		"debug",
		fmt.Sprintf("deleting terraform state for ref %s", ref),
		nil,
	)

	if !c.canLock(w, r, ref, id) {
		return
	}

	if err := c.store.DeleteState(ref); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("error deleting terraform state for ref %s", ref),
			err,
		)
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// HandleKeepVersions
func (c *Backend) HandleKeepVersions(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref in HandleKeepVersions: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// HandleListStates
func (c *Backend) HandleListStates(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref in HandleListStates: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	states, err := c.store.GetStates(ref)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to retrieve list of states for ref [%s]: %v", ref, err),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, err := json.Marshal(states)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to marshal states(%d) for ref %s: %v", len(states), ref, err),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, _ = w.Write(data)
}

// HandleListVersions
func (c *Backend) HandleListVersions(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref in HandleListVersions: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	// Check if ref came in properly
	if r.URL.Query().Get("ref") == "" {
		c.options.Logger(
			"error",
			"expecting ref as query parameter",
			err,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	versions, err := c.store.List(ref)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to retrieve list of versions for ref [%s]: %v", ref, err),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
	data, err := json.Marshal(versions)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to marshal list(%d) for ref %s: %v", len(versions), ref, err),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
	}
	_, _ = w.Write(data)
}

// HandleRetrieveVersion
func (c *Backend) HandleRetrieveVersion(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref in HandleRestoreVersion: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var versionRequest struct {
		Version string `json:"version"`
	}
	if err := json.NewDecoder(r.Body).Decode(&versionRequest); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to read version in body for ref [%s]: %v", ref, err),
			err,
		)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// get the state
	state, encrypted, err := c.store.GetState(ref, versionRequest.Version)
	if err != nil {
		if err == store.ErrNotFound {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get terraform state for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// decrypt
	if encrypted {
		decryptedState, err := c.decryptState(state)
		if err != nil {
			c.options.Logger(
				"error",
				fmt.Sprintf("failed decrypt terraform state for ref [%s]: %v", ref, err),
				err,
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		state = decryptedState
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(state)
}

// HandleRestoreVersion
func (c *Backend) HandleRestoreVersion(w http.ResponseWriter, r *http.Request) {
	ref, err := c.getRef(r)
	if err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to get ref in HandleRestoreVersion: %v", err),
			err,
		)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if err := c.Init(); err != nil {
		c.options.Logger(
			"error",
			fmt.Sprintf("failed to initialize terraform state backend for ref: %s", ref),
			err,
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

// simple interface
func toInterface(input, output interface{}) error {
	j, err := json.Marshal(input)
	if err != nil {
		return err
	}
	return json.Unmarshal(j, output)
}
