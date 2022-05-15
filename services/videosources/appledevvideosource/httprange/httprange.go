package httprange

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
)

// Error returned when the number of retries has been reached.
var ErrTooManyRetries = errors.New("Too many retries")

// ErrBadResponseCode is an error returned if the response code is non 200 to 299.
type ErrBadResponseCode int

const defaultCopyBufferSize = 4 * 1024

func (e ErrBadResponseCode) Error() string {
	return fmt.Sprintf("Non-200 HTTP response code: %d", int(e))
}

// DefaultMaxRetries is the default number of maximum retries before the
// HttpRange will give up with ErrTooManyRetries.
const DefaultMaxRetries = 30

// HTTPRange manages the configuration for fetching a range.
type HTTPRange struct {
	// requestFn is a function which returns a new request.  This request
	// will be modified to include the byte range if required.
	requestFn func() (*http.Request, error)

	// writeObserver is a function which, when set, will periodically be
	// called with updates on how much data has been written.
	writeObserver func(writtenSoFar int64, totalExpectedSize int64)

	// Maximum number of retries before giving up.
	maxRetry int
}

// new will create a new httpRange with the default configuration
func newHTTPRange() *HTTPRange {
	return &HTTPRange{
		maxRetry: DefaultMaxRetries,
	}
}

// Get will return a new HttpRange which will call the given URL with the GET request.
func Get(url string) *HTTPRange {
	httpRange := newHTTPRange()
	httpRange.requestFn = func() (*http.Request, error) {
		return http.NewRequest(http.MethodGet, url, nil)
	}
	return httpRange
}

// Sets the write observer for this HTTP Range.  This returns the original range.
func (httpRange *HTTPRange) WithWriteObserver(observer func(writtenSoFar int64, totalExpectedSize int64)) *HTTPRange {
	httpRange.writeObserver = observer
	return httpRange
}

// WriteTo will execute the request and write the response to the writer.  This will check
// that the entire response was returned, as determined by the Content-Length of the
// first request, and will retry the request from the stopping point until the response
// has been returned.
func (httpRange *HTTPRange) WriteTo(ctx context.Context, w io.Writer) (int64, error) {
	var consumedSoFar int64
	var expectedSize int64 = -1

	// Check that the request function is non-nil
	if httpRange.requestFn == nil {
		return 0, errors.New("Request function is nil")
	}

	for retry := 0; retry < httpRange.maxRetry; retry++ {
		// Sets up the request, applying a range if necessary
		req, err := httpRange.requestFn()
		if err != nil {
			return consumedSoFar, err
		}
		if (expectedSize != -1) && (consumedSoFar > 0) {
			req.Header.Set("Range", fmt.Sprintf("bytes=%d-", consumedSoFar))
		}

		req = req.WithContext(ctx)

		// Execute the request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return consumedSoFar, err
		} else if (resp.StatusCode < 200) || (resp.StatusCode > 299) {
			resp.Body.Close()
			return consumedSoFar, ErrBadResponseCode(resp.StatusCode)
		} else if ctx.Err() != nil {
			return consumedSoFar, ctx.Err()
		}

		// Determine the expected content length
		if expectedSize == -1 {
			expectedSize = resp.ContentLength
		}

		cnt, err := httpRange.doCopy(w, resp.Body, consumedSoFar, expectedSize)
		resp.Body.Close()

		consumedSoFar += cnt

		// Determine whether there is a need to continue consuming the response.
		if (err != nil) && (err != io.EOF) {
			// TEMP
			if err != io.ErrUnexpectedEOF {
				return consumedSoFar, fmt.Errorf("Cannot doCopy: %v", err)
			}
		}
		if consumedSoFar >= expectedSize {
			return consumedSoFar, nil
		}
	}

	return consumedSoFar, ErrTooManyRetries
}

// Performs a copy.  If there is a watcher, this will periodically notify the watcher.
// This returns the bytes consumed, which DOES NOT include consumedSoFar.
func (httpRange *HTTPRange) doCopy(w io.Writer, r io.Reader, consumedSoFar, expectedSize int64) (int64, error) {
	if httpRange.writeObserver == nil {
		// If there is no write observer, simply use io.Copy()
		return io.Copy(w, r)
	}

	// Otherwise, copy using a buffer manually.
	var amountCopied, written int64
	var err error
	for written, err = io.CopyN(w, r, defaultCopyBufferSize); err == nil; written, err = io.CopyN(w, r, defaultCopyBufferSize) {
		amountCopied += written
		httpRange.writeObserver(amountCopied+consumedSoFar, expectedSize)
	}
	amountCopied += written
	httpRange.writeObserver(amountCopied+consumedSoFar, expectedSize)

	return amountCopied, err
}
