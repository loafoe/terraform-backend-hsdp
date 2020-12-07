package s3

import (
	"bytes"
	"context"
	"encoding/json"

	"github.com/minio/minio-go/v7"
	"github.com/philips-labs/terraform-backend-http/backend/store"
	"github.com/philips-labs/terraform-backend-http/backend/types"
)

// GetLock gets the lock
func (c *Store) GetLock(ref string) (*types.Lock, error) {
	opts := minio.GetObjectOptions{}
	ctx := context.Background()
	lockPath := c.lockPath(ref)

	// Check if object exists
	_, err := c.client.StatObject(ctx, c.bucket, lockPath, opts)
	if err != nil {
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == "NoSuchKey" {
			return nil, store.ErrNotFound
		}
		return nil, err
	}
	object, err := c.client.GetObject(ctx, c.bucket, lockPath, opts)
	if err != nil {
		return nil, err
	}
	defer object.Close()

	var lock types.LockDocument
	if err := json.NewDecoder(object).Decode(&lock); err != nil {
		return nil, err
	}
	return &lock.Lock, nil
}

// PutLock puts the lock
func (c *Store) PutLock(ref string, lock types.Lock) error {
	lockPath := c.lockPath(ref)
	ctx := context.Background()

	document := types.LockDocument{
		Ref:  ref,
		Lock: lock,
	}
	jsonBody, err := json.Marshal(&document)
	if err != nil {
		return err
	}
	data := bytes.NewBuffer(jsonBody)
	_, err = c.client.PutObject(ctx, c.bucket, lockPath, data, int64(len(jsonBody)), minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	return nil
}

// DeleteLock deletes a lock
func (c *Store) DeleteLock(ref string) error {
	lockPath := c.lockPath(ref)
	ctx := context.Background()

	err := c.client.RemoveObject(ctx, c.bucket, lockPath, minio.RemoveObjectOptions{})

	return err
}
