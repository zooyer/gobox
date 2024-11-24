package cat

import (
	"context"
	"io"
)

type ctxReader struct {
	ctx    context.Context
	reader io.Reader
}

func (cr *ctxReader) Read(p []byte) (n int, err error) {
	select {
	case <-cr.ctx.Done():
		return 0, cr.ctx.Err()
	default:
		return cr.reader.Read(p)
	}
}

func withContext(ctx context.Context, reader io.Reader) io.Reader {
	return &ctxReader{ctx: ctx, reader: reader}
}
