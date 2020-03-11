package file

import (
	"io"
	"net/url"
	"time"

	"github.com/archy-bold/go-file-helper/schema"
	"github.com/go-kit/kit/log"
)

type loggingMiddleware struct {
	logger log.Logger
	next   schema.Service
}

func (mw *loggingMiddleware) Store(fname string, r io.Reader, size int64) (fi schema.FileInfo, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "files.Store",
			"fname", fname,
			"size", fi.Size,
			"mimetype", fi.MimeType,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	fi, err = mw.next.Store(fname, r, size)
	return
}

func (mw *loggingMiddleware) Delete(fname string) (err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "files.Delete",
			"fname", fname,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	err = mw.next.Delete(fname)
	return
}

func (mw *loggingMiddleware) GetPublicURL(fname string) (url *url.URL, err error) {
	defer func(begin time.Time) {
		mw.logger.Log(
			"method", "files.GetPublicURL",
			"fname", fname,
			"url", url,
			"err", err,
			"took", time.Since(begin),
		)
	}(time.Now())

	url, err = mw.next.GetPublicURL(fname)
	return
}
