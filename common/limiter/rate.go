package limiter

import (
	"context"
	"io"

	"github.com/xtls/xray-core/common"
	"github.com/xtls/xray-core/common/buf"
	"golang.org/x/time/rate"
)

type Writer struct {
	writer  buf.Writer
	limiter *rate.Limiter
	w       io.Writer
}

func (l *Limiter) RateWriter(writer buf.Writer, limiter *rate.Limiter) buf.Writer {
	return &Writer{
		writer:  writer,
		limiter: limiter,
	}
}

func (w *Writer) Close() error {
	return common.Close(w.writer)
}

func (w *Writer) WriteMultiBuffer(mb buf.MultiBuffer) error {
	ctx := context.Background()
	s := int(mb.Len())
	if s > w.limiter.Burst() {
		newError("mb len > burst").AtInfo().WriteToLog()
		for {
			if s <= 0 {
				newError("mb len > burst").AtInfo().WriteToLog()
				break
			}
			s -= w.limiter.Burst()
			w.limiter.WaitN(ctx, w.limiter.Burst())
		}
	}
	w.limiter.WaitN(ctx, int(mb.Len()))
	return w.writer.WriteMultiBuffer(mb)
}
