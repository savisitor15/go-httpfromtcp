package response

import (
	"fmt"
	"io"

	"github.com/savisitor15/go-httpfromtcp/internal/headers"
)

type writerState int

const (
	WriterStateStatusLine writerState = iota
	WriterStateHeaders
	WriterStateBody
	WriterStateTrailers
)

type Writer struct {
	writerState int
	writer      io.Writer
}

func WriterStatusToString(stat int) string {
	switch stat {
	case int(WriterStateStatusLine):
		return "StatusLine"
	case int(WriterStateHeaders):
		return "Headers"
	case int(WriterStateBody):
		return "Body"
	case int(WriterStateTrailers):
		return "Trailers"
	default:
		return "Unknown state"
	}
}

func (w *Writer) statusString() string {
	return WriterStatusToString(w.writerState)
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writerState: int(WriterStateStatusLine),
		writer:      w,
	}
}

func (w *Writer) WriteStatusLine(statusCode StatusCode) error {
	if w.writerState != int(WriterStateStatusLine) {
		return fmt.Errorf("error invalid status: %d - %s", w.writerState, w.statusString())
	}
	defer func() { w.writerState = int(WriterStateHeaders) }()
	_, err := w.writer.Write(getStatusLine(statusCode))
	return err
}

func (w *Writer) WriteHeaders(headers headers.Headers) error {
	if w.writerState != int(WriterStateHeaders) {
		return fmt.Errorf("error cannot write headers in current state: %d - %s", w.writerState, w.statusString())
	}
	defer func() { w.writerState = int(WriterStateBody) }()
	for k, v := range headers {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}

func (w *Writer) WriteBody(p []byte) (int, error) {
	if w.writerState != int(WriterStateBody) {
		return 0, fmt.Errorf("error cannot write body in current state: %d - %s", w.writerState, w.statusString())
	}
	return w.writer.Write(p)
}

func (w *Writer) WriteChunkedBody(p []byte) (int, error) {
	lineEnd := []byte("\r\n")
	if w.writerState != int(WriterStateBody) {
		return 0, fmt.Errorf("error cannot write body in current state: %d - %s", w.writerState, w.statusString())
	}
	leng := len(p)
	// write length indicator
	blen, err := w.WriteBody([]byte(fmt.Sprintf("%x\r\n", leng)))
	if err != nil {
		return 0, err
	}
	payload := append(p, lineEnd...)
	b2len, err := w.WriteBody(payload)
	if err != nil {
		return blen, err
	}
	return blen + b2len, err
}

func (w *Writer) WriteChunkedBodyDone() (int, error) {
	if w.writerState != int(WriterStateBody) {
		return 0, fmt.Errorf("error cannot write body in current state: %d - %s", w.writerState, w.statusString())
	}
	defer func() { w.writerState = int(WriterStateTrailers) }()
	return w.WriteBody([]byte("0\r\n"))
}

func (w *Writer) WriteTrailers(h headers.Headers) error {
	if w.writerState != int(WriterStateTrailers) {
		return fmt.Errorf("error cannot write trailers in current state: %d - %s", w.writerState, w.statusString())
	}
	for k, v := range h {
		_, err := w.writer.Write([]byte(fmt.Sprintf("%s: %s\r\n", k, v)))
		if err != nil {
			return err
		}
	}
	_, err := w.writer.Write([]byte("\r\n"))
	return err
}
