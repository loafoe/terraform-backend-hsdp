package s3

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/loafoe/terraform-backend-hsdp/backend/store"
)

var _ store.Stats = (*Store)(nil)

// Locks gets the lock
func (c *Store) Locks(age int) (int, error) {
	ctx := context.Background()

	lockPath := c.lockPath("") // Base

	ageTimestamp := time.Now().Add(-time.Second * time.Duration(86400*age))

	opts := minio.ListObjectsOptions{
		Prefix:       lockPath,
		Recursive:    true,
		WithMetadata: true,
	}
	count := 0
	ch := c.client.ListObjects(ctx, c.bucket, opts)
	for object := range ch {
		if object.Err != nil {
			fmt.Println(object.Err)
			continue
		}
		// Count only older locks
		if age > 0 {
			if object.LastModified.Before(ageTimestamp) {
				count = count + 1
			}
			continue
		}
		count = count + 1
	}
	return count, nil
}

// States gets the lock
func (c *Store) States() (int, error) {
	ctx := context.Background()

	storePath := c.storePath("") // Base

	opts := minio.ListObjectsOptions{
		Prefix:       storePath,
		Recursive:    true,
		WithMetadata: true,
	}
	count := 0
	ch := c.client.ListObjects(ctx, c.bucket, opts)
	for object := range ch {
		if object.Err != nil {
			fmt.Println(object.Err)
			continue
		}
		count = count + 1
	}
	return count, nil
}

// Identities gets the lock
func (c *Store) Identities() (int, error) {
	ctx := context.Background()

	storePath := c.storePath("") // Base

	opts := minio.ListObjectsOptions{
		Prefix:       storePath,
		Recursive:    true,
		WithMetadata: true,
	}
	ch := c.client.ListObjects(ctx, c.bucket, opts)

	var ids []string
	for object := range ch {
		if object.Err != nil {
			fmt.Println(object.Err)
			continue
		}
		parts := strings.Split(object.Key, "/")
		if len(parts) > 2 {
			ids = append(ids, parts[2])
		}
	}

	return len(unique(ids)), nil
}

func unique(stringSlice []string) []string {
	keys := make(map[string]bool)
	list := []string{}
	for _, entry := range stringSlice {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}
