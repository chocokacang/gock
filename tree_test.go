package gock

import (
	"testing"
)

func TestCountSections(t *testing.T) {
	cases := []struct {
		name     string
		path     string
		expected uint16
	}{
		{name: "empty string", path: "", expected: 0},
		{name: "root path", path: "/", expected: 1},
		{name: "single section", path: "/foo", expected: 1},
		{name: "multiple sections", path: "/foo/bar/baz", expected: 3},
		{name: "trailing slash", path: "/foo/bar/", expected: 3},
		{name: "consecutive slashes", path: "///foo/bar", expected: 4},
		{name: "non-slash path", path: "foo/bar", expected: 1},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := countSections(tc.path)
			if actual != tc.expected {
				t.Errorf("countSections(%q) = %d, expected %d", tc.path, actual, tc.expected)
			}
		})
	}
}
