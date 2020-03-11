package file

import (
	"bytes"
	"testing"

	"github.com/archy-bold/go-file-helper/s3"
	"github.com/archy-bold/go-file-helper/schema"
	"github.com/go-kit/kit/log"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	minioCreds = schema.Credentials{
		"play.minio.io:9000",
		"Q3AM3UQ867SPQQA43P2F",
		"zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		true,
		"mymusic",
		"us-east-1",
	}
)

var b bytes.Buffer

var newServiceTests = map[string]struct {
	typ      string
	creds    schema.Credentials
	expected schema.Service
	logger   log.Logger
	err      error
}{
	"s3 service":          {"s3", minioCreds, &s3.Service{}, nil, nil},
	"s3 with logger":      {"s3", minioCreds, &loggingMiddleware{}, log.NewLogfmtLogger(&b), nil},
	"unsupported service": {"disk", schema.Credentials{}, nil, nil, errors.Wrap(ErrUnsupportedType, "service 'disk' unsupported")},
}

func Test_NewService(t *testing.T) {
	for tn, tt := range newServiceTests {
		svc, err := NewService(tt.typ, tt.creds, tt.logger)

		if tt.err == nil {
			assert.NoErrorf(t, err, "Expected no err on NewService test '%s'", tn)
			assert.IsTypef(t, tt.expected, svc, "Expected service type to match on NewService test '%s'", tn)
		} else {
			assert.Nilf(t, svc, "Expected svc to be nil on NewService test '%s'", tn)
			assert.EqualErrorf(t, err, tt.err.Error(), "Expected err to match on NewService test '%s'", tn)
		}
	}
}
