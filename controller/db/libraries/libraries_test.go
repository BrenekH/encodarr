package libraries

import (
	"fmt"
	"reflect"
	"testing"
	"time"
)

func TestFromDBLibrary(t *testing.T) {
	tests := []struct {
		name   string
		errors bool
		in     dBLibrary
		want   Library
	}{
		{name: "All parameters (no error)", errors: false, in: dBLibrary{
			FsCheckInterval: "1h1m1s",
			Pipeline:        []byte("{}"),
			Queue:           []byte("{}"),
			FileCache:       []byte("{}"),
			PathMasks:       []byte("[]"),
		}, want: Library{
			FsCheckInterval: time.Duration(time.Hour + time.Minute + time.Second),
			Pipeline:        pluginPipeline{},
			Queue:           queue{},
			FileCache:       fileCache{},
			PathMasks:       []string{},
		}},
		// Error tests
		{name: "FsCheckInterval causes error", errors: true, in: dBLibrary{FsCheckInterval: "Invalid"}},
		{name: "Pipeline causes error", errors: true, in: dBLibrary{Pipeline: []byte("Invalid")}},
		{name: "Queue causes error", errors: true, in: dBLibrary{Queue: []byte("Invalid")}},
		{name: "FileCache causes error", errors: true, in: dBLibrary{FileCache: []byte("Invalid")}},
		{name: "PathMasks causes error", errors: true, in: dBLibrary{PathMasks: []byte("Invalid")}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.name)

		t.Run(testname, func(t *testing.T) {
			l := Library{}

			ansErr := l.fromDBLibrary(tt.in)

			if tt.errors && ansErr == nil {
				t.Errorf("expected an error from %v", tt.in)
			} else if tt.errors && ansErr != nil {
				return
			}

			if !reflect.DeepEqual(l, tt.want) {
				t.Errorf("got %v, want %v", l, tt.want)
			}
		})
	}
}

func TestToDBLibrary(t *testing.T) {
	tests := []struct {
		name   string
		errors bool
		in     Library
		want   dBLibrary
	}{
		{name: "Basic", errors: false, in: Library{
			FsCheckInterval: time.Duration(time.Hour + time.Minute + time.Second),
			Pipeline:        pluginPipeline{},
			Queue:           queue{},
			FileCache:       fileCache{},
			PathMasks:       []string{},
		}, want: dBLibrary{
			FsCheckInterval: "1h1m1s",
			Pipeline:        []byte("{}"),
			Queue:           []byte("{}"),
			FileCache:       []byte("{}"),
			PathMasks:       []byte("[]"),
		}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.name)

		t.Run(testname, func(t *testing.T) {
			ans, ansErr := tt.in.toDBLibrary()

			if tt.errors && ansErr == nil {
				t.Errorf("expected an error from %v", tt.in)
			} else if tt.errors && ansErr != nil {
				return
			}

			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
