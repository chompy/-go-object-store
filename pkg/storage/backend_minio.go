package storage

import (
	"bytes"
	"context"
	"io/ioutil"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"gitlab.com/contextualcode/storage-backend/pkg/types"

	"github.com/pkg/errors"
)

// MinioBackend handles connection to S3 compatible backend.
type MinioBackend struct {
	Endpoint        string
	AccessKeyID     string
	SecretAccessKey string
	Bucket          string
	client          *minio.Client
}

func (b *MinioBackend) init() error {
	if b.client != nil {
		return nil
	}
	var err error
	b.client, err = minio.New(b.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(b.AccessKeyID, b.SecretAccessKey, ""),
		Secure: true,
	})
	return errors.WithStack(err)
}

// Upload uploads given object.
func (b *MinioBackend) Upload(o *types.Object) error {
	if err := b.init(); err != nil {
		return errors.WithStack(err)
	}
	data, err := o.Serialize()
	if err != nil {
		return errors.WithStack(err)
	}
	if _, err := b.client.PutObject(
		context.Background(),
		b.Bucket,
		o.UID,
		bytes.NewReader(data),
		int64(len(data)),
		minio.PutObjectOptions{},
	); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// Download downloads object.
func (b *MinioBackend) Download(uid string) (*types.Object, error) {
	if err := b.init(); err != nil {
		return nil, errors.WithStack(err)
	}
	data, err := b.client.GetObject(
		context.Background(),
		b.Bucket,
		uid,
		minio.GetObjectOptions{},
	)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	o := &types.Object{}
	dataBytes, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if err := o.Unserialize(dataBytes); err != nil {
		return nil, errors.WithStack(err)
	}
	return o, nil
}
