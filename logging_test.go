package file

import (
	"bytes"
	"fmt"
	"net/url"
	"testing"

	mocks "github.com/archy-bold/go-file-helper/mocks/schema"
	"github.com/archy-bold/go-file-helper/schema"
	"github.com/go-kit/kit/log"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	exampleURL, _ = url.Parse("http://localhost/example.xml")
)

var storeLoggingTests = map[string]struct {
	fname string
	fi    schema.FileInfo
	err   error
}{
	"success": {"test.xml", schema.FileInfo{Size: 1000, MimeType: "application/xml"}, nil},
	"error":   {"error.txt", schema.FileInfo{}, assert.AnError},
}

func Test_logging_Store(t *testing.T) {
	r := bytes.NewBuffer(make([]byte, 0))

	for tn, tt := range storeLoggingTests {
		// Mock the service
		var b bytes.Buffer
		mSvc, svc := mockService(&b)
		mSvc.On("Store", tt.fname, r, int64(0)).
			Return(tt.fi, tt.err).
			Once()

		// Run the function
		_, err := svc.Store(tt.fname, r, int64(0))

		// Assert it's as expected
		mSvc.AssertCalled(t, "Store", tt.fname, r, int64(0))
		expectedErr := "null"
		if err != nil {
			expectedErr = err.Error()
		}
		logs := []string{
			"method=files.Store",
			fmt.Sprintf("fname=%s", tt.fname),
			fmt.Sprintf("size=%d", tt.fi.Size),
			fmt.Sprintf("mimetype=%s", tt.fi.MimeType),
			fmt.Sprintf("err=%s", expectedErr),
		}
		assertLogging(t, b, 6, logs, tn)
	}
}

var deleteLoggingTests = map[string]struct {
	fname string
	err   error
}{
	"success": {"test.xml", nil},
	"error":   {"error.txt", assert.AnError},
}

func Test_logging_Delete(t *testing.T) {
	for tn, tt := range deleteLoggingTests {
		// Mock the service
		var b bytes.Buffer
		mSvc, svc := mockService(&b)
		mSvc.On("Delete", tt.fname).
			Return(tt.err).
			Once()

		// Run the function
		svc.Delete(tt.fname)

		// Assert it's as expected
		mSvc.AssertCalled(t, "Delete", tt.fname)
		expectedErr := "null"
		if tt.err != nil {
			expectedErr = tt.err.Error()
		}
		logs := []string{
			"method=files.Delete",
			fmt.Sprintf("fname=%s", tt.fname),
			fmt.Sprintf("err=%s", expectedErr),
		}
		assertLogging(t, b, 4, logs, tn)
	}
}

var getPublicURLLoggingTests = map[string]struct {
	fname string
	url   *url.URL
	err   error
}{
	"success": {"test.xml", exampleURL, nil},
	"error":   {"error.txt", nil, assert.AnError},
}

func Test_logging_GetPublicURL(t *testing.T) {
	for tn, tt := range getPublicURLLoggingTests {
		// Mock the service
		var b bytes.Buffer
		mSvc, svc := mockService(&b)
		mSvc.On("GetPublicURL", tt.fname).
			Return(tt.url, tt.err).
			Once()

		// Run the function
		svc.GetPublicURL(tt.fname)

		// Assert it's as expected
		mSvc.AssertCalled(t, "GetPublicURL", tt.fname)
		expectedErr := "null"
		expectedURL := "null"
		if tt.err != nil {
			expectedErr = tt.err.Error()
		}
		if tt.url != nil {
			expectedURL = tt.url.String()
		}
		logs := []string{
			"method=files.GetPublicURL",
			fmt.Sprintf("fname=%s", tt.fname),
			fmt.Sprintf("url=%s", expectedURL),
			fmt.Sprintf("err=%s", expectedErr),
		}
		assertLogging(t, b, 5, logs, tn)
	}
}

func mockService(b *bytes.Buffer) (*mocks.Service, schema.Service) {
	mSvc := &mocks.Service{}
	logger := log.NewLogfmtLogger(b)
	var svc schema.Service
	svc = &loggingMiddleware{logger, mSvc}
	return mSvc, svc
}

func assertLogging(t *testing.T, b bytes.Buffer, len int, ss []string, test string) {
	// Assert we've logged the call
	logs := b.String()
	logging, err := shellwords.Parse(logs)
	require.Nil(t, err)
	assert.Lenf(t, logging, len, "Expected different logging length for test '%s'", test)
	for i, s := range ss {
		assert.Equal(t, s, logging[i], "Expected different logs for test '%s'", test)
	}
}
