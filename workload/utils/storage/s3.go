package storage

import (
	"strings"

	"github.com/minio/minio-go"
	"github.com/pingcap/br/pkg/storage"
	"github.com/pingcap/errors"
	"github.com/pingcap/kvproto/pkg/backup"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

// TempS3Storage is a temporary S3 storage.
// currently, this is only use for storing BR backup result.
type TempS3Storage struct {
	opts  *backup.S3
	Raw   string
	minio *minio.Client
}

// ConnectToS3 creates a new TempS3Storage by the url string.
func ConnectToS3(url string) (*TempS3Storage, error) {
	backend, err := storage.ParseBackend(url, &storage.BackendOptions{})
	if err != nil {
		return nil, err
	}
	s3Opts := backend.GetS3()
	log.Info("use temporary S3 storage", zap.Any("config", s3Opts), zap.String("url", url))
	minioClient, err := minio.New(strings.TrimPrefix(s3Opts.Endpoint, "http://"), s3Opts.AccessKey, s3Opts.SecretAccessKey, false)
	if err != nil {
		return nil, err
	}
	return &TempS3Storage{opts: s3Opts, Raw: url, minio: minioClient}, nil
}

// Cleanup drops contents of this storage
func (s *TempS3Storage) Cleanup() error {
	for obj := range s.minio.ListObjects(s.opts.Bucket, s.opts.Prefix, true, nil) {
		if obj.Err != nil {
			return errors.Trace(obj.Err)
		}
		// TODO batching the request
		if err := s.minio.RemoveObject(s.opts.Bucket, obj.Key); err != nil {
			return err
		}
	}
	return nil
}
