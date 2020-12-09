package s3

import (
	"bytes"
	"context"
	"encoding/json"
	"path/filepath"

	"github.com/minio/minio-go/v7"

	"github.com/philips-labs/terraform-backend-hsdp/backend/store"
	"github.com/philips-labs/terraform-backend-hsdp/backend/types"
)

func (c *Store) storePath(ref string) string {
	return filepath.Join("store", ref)
}

func (c *Store) lockPath(ref string) string {
	return filepath.Join("lock", ref)
}

// GetState gets the state
func (c *Store) GetState(ref string) (map[string]interface{}, bool, error) {
	opts := minio.GetObjectOptions{}
	storePath := c.storePath(ref)
	ctx := context.Background()

	// Check if object exists
	_, err := c.client.StatObject(ctx, c.bucket, storePath, opts)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return nil, false, store.ErrNotFound
		}
		return nil, false, err
	}
	object, err := c.client.GetObject(ctx, c.bucket, storePath, opts)
	if err != nil {
		return nil, false, err
	}
	defer object.Close()

	var state types.StateDocument
	if err := json.NewDecoder(object).Decode(&state); err != nil {
		return nil, false, err
	}
	return state.State, state.Encrypted, nil
}

// PutState puts the state
func (c *Store) PutState(ref string, state, metadata map[string]interface{}, encrypted bool) error {
	storePath := c.storePath(ref)
	ctx := context.Background()

	document := types.StateDocument{
		Ref:       ref,
		State:     state,
		Encrypted: encrypted,
		Metadata:  metadata,
	}
	jsonBody, err := json.Marshal(&document)
	if err != nil {
		return err
	}
	data := bytes.NewBuffer(jsonBody)

	_, err = c.client.PutObject(ctx, c.bucket, storePath, data, int64(len(jsonBody)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DeleteState deletes a state
func (c *Store) DeleteState(ref string) error {
	storePath := c.storePath(ref)
	ctx := context.Background()

	err := c.client.RemoveObject(ctx, c.bucket, storePath, minio.RemoveObjectOptions{})

	return err
}
