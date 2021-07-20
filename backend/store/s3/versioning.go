package s3

import (
	"context"
	"fmt"
	"path"

	"github.com/minio/minio-go/v7"
)

func (c *Store) List(ref string) ([]string, error) {
	var versions []string
	versionFolder := c.versionFolder(ref)
	ctx := context.Background()

	opts := minio.ListObjectsOptions{
		Prefix:    versionFolder,
		Recursive: true,
	}
	ch := c.client.ListObjects(ctx, c.bucket, opts)
	for object := range ch {
		if object.Err != nil {
			fmt.Println(object.Err)
			continue
		}
		_, key := path.Split(object.Key)
		versions = append(versions, key)
	}
	return versions, nil
}

func (c *Store) Keep(last int) error {
	return fmt.Errorf("implement me")
}

func (c *Store) Restore(ref, version string) error {
	return fmt.Errorf("implement me")
}
