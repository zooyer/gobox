package cat

import (
	"bytes"
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"
)

func TestWithContext(t *testing.T) {
	tests := []struct {
		name         string
		ctxFunc      func() (context.Context, context.CancelFunc)
		reader       io.Reader
		expectedData string
		expectedErr  error
	}{
		{
			name: "Read with no cancellation",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			reader:       bytes.NewReader([]byte("Hello, world!")),
			expectedData: "Hello, world!",
			expectedErr:  nil,
		},
		{
			name: "Read with cancellation",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				defer cancel()
				return ctx, cancel
			},
			reader:       strings.NewReader("Hello, world!"),
			expectedData: "",
			expectedErr:  context.Canceled,
		},
		{
			name: "Read with timeout",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)

				time.Sleep(time.Millisecond * 50)

				return ctx, cancel
			},
			reader:       strings.NewReader("Hello, world!"),
			expectedData: "",
			expectedErr:  context.DeadlineExceeded,
		},
		{
			name: "EOF handling",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				return context.WithCancel(context.Background())
			},
			reader:       bytes.NewReader([]byte("")),
			expectedData: "",
			expectedErr:  io.EOF,
		},
		{
			name: "Cancel after reading some data",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				go func() {
					time.Sleep(50 * time.Millisecond)
					cancel()
				}()
				return ctx, cancel
			},
			reader:       strings.NewReader("Hello, world!"),
			expectedData: "Hello, world!",
			expectedErr:  context.Canceled,
		},
		{
			name: "Read with immediate cancel",
			ctxFunc: func() (context.Context, context.CancelFunc) {
				ctx, cancel := context.WithCancel(context.Background())
				cancel() // Cancel immediately
				return ctx, cancel
			},
			reader:       strings.NewReader("Hello, world!"),
			expectedData: "",
			expectedErr:  context.Canceled,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := tt.ctxFunc()
			defer cancel()

			reader := withContext(ctx, tt.reader)
			buf := make([]byte, 1024)
			var result string

			for {
				n, err := reader.Read(buf)
				result += string(buf[:n])

				if err != nil {
					if !errors.Is(err, tt.expectedErr) && err != io.EOF {
						t.Fatalf("Expected error %v, got %v", tt.expectedErr, err)
					}
					break
				}
			}

			if result != tt.expectedData {
				t.Fatalf("Expected data %q, got %q", tt.expectedData, result)
			}
		})
	}
}
