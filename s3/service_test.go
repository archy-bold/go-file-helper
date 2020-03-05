package s3

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/archy-bold/go-file-helper/schema"
	minio "github.com/minio/minio-go"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var (
	testfile, _       = os.Open("testdata/test.xml")
	b                 = []byte{}
	emptyreader       = bytes.NewReader(b)
	maxSize     int64 = (1024 * 1024 * 1024 * 1024 * 5) + 1
)

var storeTests = map[string]struct {
	fname    string
	reader   io.Reader
	size     int64
	expected schema.FileInfo
	err      error
}{
	"success": {"test/test.xml", testfile, 20, schema.FileInfo{"test/test.xml", "mymusic", "us-east-1", "application/xml", 20, time.Now()}, nil},
	"error":   {"test/empty.xml", emptyreader, maxSize, schema.FileInfo{}, errors.Wrap(errors.New("Your proposed upload size ‘5497558138881’ exceeds the maximum allowed object size ‘5497558138880’ for single PUT operation."), "failed to store file 'test/empty.xml'")},
}

func Test_Service_Store(t *testing.T) {
	svc := getServiceForTesting()

	for tn, tt := range storeTests {
		now := time.Now()
		ret, err := svc.Store(tt.fname, tt.reader, tt.size)

		if tt.err == nil {
			assert.NoErrorf(t, err, "Expected no err on Store test '%s'", tn)
			assertFileInfoMatchf(t, tt.expected, ret, now, "Expected FileInfo to match on Store test '%s'", tn)
		} else {
			assert.EqualErrorf(t, err, tt.err.Error(), "Expected err to match on Store test '%s'", tn)
			assert.Equalf(t, ret, tt.expected, "Expected FileInfo to match on Store test '%s'", tn)
		}
	}
}

var getPublicURLTests = map[string]struct {
	fname string
	err   error
}{
	"success": {"test/test.xml", nil},
	"error":   {"", errors.Wrap(errors.New("Object name cannot be empty"), "unable to get S3 pre-signed URL for file ''")},
}

func Test_Service_GetPublicURL(t *testing.T) {
	svc := getServiceForTesting()

	for tn, tt := range getPublicURLTests {
		url, err := svc.GetPublicURL(tt.fname)

		if tt.err == nil {
			assert.NoErrorf(t, err, "Expected no err on GetPublicURL test '%s'", tn)
			assert.Truef(t, strings.HasPrefix(url.String(), "https://play.minio.io:9000/mymusic"), "Expected URL to be correct on GetPublicURL test '%s'", tn)
		} else {
			assert.EqualErrorf(t, err, tt.err.Error(), "Expected err to match on GetPublicURL test '%s'", tn)
			assert.Nilf(t, url, "Expected url to be nil on GetPublicURL test '%s'", tn)
		}
	}
}

var deleteTests = map[string]struct {
	fname string
	err   error
}{
	"success": {"test/test.xml", nil},
	"error":   {"", errors.Wrap(errors.New("Object name cannot be empty"), "cannot delete ''")},
}

func Test_Service_Delete(t *testing.T) {
	svc := getServiceForTesting()

	for tn, tt := range deleteTests {
		err := svc.Delete(tt.fname)

		if tt.err == nil {
			assert.NoErrorf(t, err, "Expected no err on Delete test '%s'", tn)
		} else {
			assert.EqualErrorf(t, err, tt.err.Error(), "Expected err to match on Delete test '%s'", tn)
		}
	}
}

var newServiceTests = map[string]struct {
	creds schema.Credentials
	err   error
}{
	"success":           {schema.Credentials{"play.minio.io:9000", "Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", true, "mymusic", ""}, nil},
	"invalid creds":     {schema.Credentials{}, errors.Wrap(minio.ErrInvalidArgument("Endpoint:  does not follow ip address or domain name standards."), "cannot create s3 client")},
	"invalid bucket":    {schema.Credentials{"play.minio.io:9000", "Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", true, "", ""}, errors.New("bucket error: Bucket name cannot be empty")},
	"bucket not exists": {schema.Credentials{"play.minio.io:9000", "Q3AM3UQ867SPQQA43P2F", "zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG", true, "bleepbloop", ""}, errors.New("bucket bleepbloop: bucket doesn't exist")},
}

func Test_NewService(t *testing.T) {
	for tn, tt := range newServiceTests {
		svc, err := NewService(tt.creds)

		if tt.err == nil {
			assert.Nilf(t, err, "Expected err to be nil on NewService test '%s'", tn)
			assert.NotNilf(t, svc, "Expected svc not to be nil on NewService test '%s'", tn)
		} else {
			assert.Nilf(t, svc, "Expected svc to be nil on NewService test '%s'", tn)
			assert.EqualErrorf(t, err, tt.err.Error(), "Expected err to match on NewService test '%s'", tn)
		}
	}
}

func assertFileInfoMatchf(t *testing.T, expected schema.FileInfo, actual schema.FileInfo, d time.Time, msg string, opts ...interface{}) {
	assert.Equalf(t, expected.FullPath, actual.FullPath, msg, opts...)
	assert.Equalf(t, expected.Bucket, actual.Bucket, msg, opts...)
	assert.Equalf(t, expected.Region, actual.Region, msg, opts...)
	assert.Equalf(t, expected.MimeType, actual.MimeType, msg, opts...)
	assert.Equalf(t, expected.Size, actual.Size, msg, opts...)
	assert.WithinDurationf(t, d, actual.CreatedAt, 5*time.Second, msg, opts...)
}

func getServiceForTesting() Service {
	creds := schema.Credentials{
		"play.minio.io:9000",
		"Q3AM3UQ867SPQQA43P2F",
		"zuf+tfteSlswRu7BJ86wekitnifILbZam1KYY3TG",
		true,
		"mymusic",
		"us-east-1",
	}

	client, _ := minio.New(creds.Endpoint, creds.ClientID, creds.ClientSecret, true)

	return Service{creds, client}
}
