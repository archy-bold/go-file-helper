package main

import (
	"github.com/archy-bold/go-file-helper/s3"
	"github.com/archy-bold/go-file-helper/schema"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
)

// ErrUnsupportedType this error indicates the service type isn't supported
var ErrUnsupportedType = errors.New("the given service type is unsupported")

// NewService returns a file service for the given type, eg S3
func NewService(typ string, creds schema.Credentials, logger log.Logger) (svc schema.Service, err error) {
	if typ == "s3" {
		svc, err = s3.NewService(creds)
	}

	if svc == nil {
		err = errors.Wrapf(ErrUnsupportedType, "service '%s' unsupported", typ)
		return
	}

	// Wrap the logger middleware
	if logger != nil {
		svc = &loggingMiddleware{logger, svc}
	}

	return
}
