package s3

import (
	"fmt"

	"github.com/cloudfoundry-community/gautocloud"
	"github.com/minio/minio-go/v7"
	"github.com/philips-software/gautocloud-connectors/hsdp"
)

// Options S3 backend options
type Options struct {
}

// NewStore creates a new S3 backend
func NewStore(opts *Options) *Store {
	if opts == nil {
		opts = &Options{}
	}
	backend := Store{}
	return &backend
}

// Store S3 store
type Store struct {
	client *minio.Client
	bucket string
}

// Init initializes the backend
func (c *Store) Init() error {
	// if there is no client then connect
	if c.client == nil {
		var svc *hsdp.S3MinioClient

		err := gautocloud.Inject(&svc)
		if err != nil {
			return fmt.Errorf("gautocloud inject: %w", err)
		}
		c.client = svc.Client
		c.bucket = svc.Bucket
	}
	return nil
}
