package main

import (
	"io"
	"testing"
	"time"
)

type MockReader struct {
	bytes []byte
}

func (m *MockReader) Read(p []byte) (n int, err error) {
	m.bytes = append(m.bytes, p...)
	return len(p), nil
}

type MockWriter struct {
	bytes []byte
}

func (m *MockWriter) Write(p []byte) (n int, err error) {
	m.bytes = append(m.bytes, p...)
	return len(p), nil
}

func Test_detectAndSendInput(t *testing.T) {
	tests := []struct {
		name   string
		r      io.Reader
		w      *MockWriter
		input  string
		output string
	}{
		{name: "scanInput1", r: &MockReader{}, w: &MockWriter{}, input: "test\n", output: "test3"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			go func() {
				detectAndSendInput(tt.r, tt.w)
				if string(tt.w.bytes) != tt.output {
					t.Errorf("%q. detectAndSendInput() = %v, want %v", tt.name, tt.input, tt.output)
				}
			}()
			time.Sleep(1000 * time.Millisecond)
			tt.r.Read([]byte(tt.input))
		})
	}
}
