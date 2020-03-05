package schema

import (
	"io"
	"net/url"
	"time"
)

// Service represents a service that can store files
type Service interface {
	Store(fname string, r io.Reader, size int64) (FileInfo, error)
	Delete(fname string) error
	GetPublicURL(fname string) (*url.URL, error)
}

// FileInfo describes a file stored on a file service
type FileInfo struct {
	FullPath  string
	Bucket    string
	Region    string
	MimeType  string
	Size      int64
	CreatedAt time.Time
}

// Credentials represents a collection of credentials that might be needed to
// connect to a file service
type Credentials struct {
	Endpoint     string
	ClientID     string
	ClientSecret string
	UseSSL       bool
	Bucket       string
	Region       string
}
