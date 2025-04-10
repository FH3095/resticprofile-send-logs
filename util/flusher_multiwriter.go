package util

import (
	"io"

	"github.com/creativeprojects/resticprofile/term"
)

type FlusherWriter interface {
	io.Writer
	term.Flusher
}

type flushableMultiWriter struct {
	writers []io.Writer
}

func (t flushableMultiWriter) Write(p []byte) (int, error) {
	for _, writer := range t.writers {
		if writtenBytes, err := writer.Write(p); err != nil {
			return writtenBytes, err
		} else if writtenBytes != len(p) {
			return writtenBytes, io.ErrShortWrite
		}
	}
	return len(p), nil
}

func (t flushableMultiWriter) Flush() error {
	for _, writer := range t.writers {
		if flusher, ok := writer.(term.Flusher); ok {
			if err := flusher.Flush(); err != nil {
				return err
			}
		}
	}
	return nil
}

func FlusherMultiWriter(writers ...io.Writer) FlusherWriter {
	return &flushableMultiWriter{writers}
}
