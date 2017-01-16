package youtube

import (
	"io"
	"math"
	"net/http"
	"sync/atomic"
)

type Reader struct {
	io.Reader
	totalSize   int64
	currentSize int64
	ch          chan<- float64
}

func (r *Reader) Read(p []byte) (n int, err error) {
	n, err = r.Reader.Read(p)
	r.ch <- math.Trunc(float64(float64(r.Add(int64(n)))/float64(r.totalSize)) * 100.0)
	return
}

// Add will count the copy current size
func (r *Reader) Add(n int64) int64 {
	return atomic.AddInt64(&r.currentSize, int64(n))
}

// Close the reader when it implements io.Closer
func (r *Reader) Close() (err error) {
	if closer, ok := r.Reader.(io.Closer); ok {
		return closer.Close()
	}
	return
}

func newReaderWithPercent(r *http.Response, per chan<- float64) *Reader {
	return &Reader{r.Body, r.ContentLength, 0, per}
}
