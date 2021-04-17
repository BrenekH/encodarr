package runner

import (
	"context"
	"testing"
)

func TestIsContextFinished(t *testing.T) {
	fCtx, cancel := context.WithCancel(context.Background())
	cancel()
	uCtx := context.Background()

	tests := []struct {
		name string
		in   context.Context
		out  bool
	}{
		{name: "Finished Context", in: fCtx, out: true},
		{name: "Unfinished Context", in: uCtx, out: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			out := IsContextFinished(&test.in)

			if out != test.out {
				t.Errorf("expected %v but got %v", test.out, out)
			}
		})
	}
}
