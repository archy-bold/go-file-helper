package s3

import (
	"fmt"
	"io"
	"mime"
	"net/url"
	"path/filepath"
	"time"

	files "github.com/archy-bold/go-file-helper/schema"
	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
)

// ErrBucketNotExists an error indicating the given bucket doesn't exist
var ErrBucketNotExists = errors.New("bucket doesn't exist")

// Service is an implementation of the files Service interface for sending
// files to S3 compliant services
type Service struct {
	creds  files.Credentials
	client *minio.Client
}

// Store store a file on the S3 service
func (s3 *Service) Store(fname string, r io.Reader, size int64) (files.FileInfo, error) {
	// Get the extension
	ext := filepath.Ext(fname)
	mtype := mime.TypeByExtension(ext)

	// Now put the object
	opts := minio.PutObjectOptions{
		ContentType: mtype,
	}
	actualSize, err := s3.client.PutObject(s3.creds.Bucket, fname, r, size, opts)

	if err != nil {
		return files.FileInfo{}, errors.Wrapf(err, "failed to store file '%s'", fname)
	}

	return files.FileInfo{
		FullPath:  fname,
		Bucket:    s3.creds.Bucket,
		Region:    s3.creds.Region,
		MimeType:  mtype,
		Size:      actualSize,
		CreatedAt: time.Now(),
	}, nil
}

// Delete deletes the given file
func (s3 *Service) Delete(fname string) (err error) {
	// Simply delete
	err = s3.client.RemoveObject(s3.creds.Bucket, fname)

	if err != nil {
		err = errors.Wrapf(err, "cannot delete '%s'", fname)
	}

	return
}

// GetPublicURL gets the public-facing URL for the given file
func (s3 *Service) GetPublicURL(fname string) (*url.URL, error) {
	// Set request parameters for content-disposition ie download.
	basename := filepath.Base(fname)
	cdHeader := fmt.Sprintf("attachment; filename=\"%s\"", basename)
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", cdHeader)

	// Generates a presigned url which expires in a day.
	url, err := s3.client.PresignedGetObject(s3.creds.Bucket, fname, time.Hour*24, reqParams)

	if err != nil {
		return nil, errors.Wrapf(err, "unable to get S3 pre-signed URL for file '%s'", fname)
	}

	return url, nil
}

// NewService gets an Service for the given credentials
func NewService(creds files.Credentials) (*Service, error) {
	// Create the client
	client, err := minio.New(creds.Endpoint, creds.ClientID, creds.ClientSecret, creds.UseSSL)

	if err != nil {
		return nil, errors.Wrap(err, "cannot create s3 client")
	}

	// Check the given bucket exists
	var exists bool
	exists, err = client.BucketExists(creds.Bucket)

	if err != nil {
		return nil, errors.Wrap(err, "bucket error")
	} else if !exists {
		return nil, errors.Wrapf(ErrBucketNotExists, "bucket %s", creds.Bucket)
	}

	return &Service{creds, client}, nil
}
