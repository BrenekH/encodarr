package controller

import (
	"fmt"
	"reflect"
	"testing"
)

func TestFilterNonVideoExts(t *testing.T) {
	var tests = []struct {
		a    []string
		want []string
	}{
		{[]string{"/input/name.txt"}, []string{}},
		{[]string{"/input/many.mp4", "/input/not this one though.notmp4"}, []string{"/input/many.mp4"}},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			ans := filterNonVideoExts(tt.a)
			if !reflect.DeepEqual(ans, tt.want) {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}

func TestIsVideoFileExt(t *testing.T) {
	var tests = []struct {
		a    string
		want bool
	}{
		{".mkv", true},
		{".mp4", true},
		{".avi", true},
		{".mka", false},
		{".mks", false},
		{".txt", false},
	}

	for _, tt := range tests {
		testname := fmt.Sprintf("%v", tt.a)
		t.Run(testname, func(t *testing.T) {
			ans := isVideoFileExt(tt.a)
			if ans != tt.want {
				t.Errorf("got %v, want %v", ans, tt.want)
			}
		})
	}
}
