package main

import "testing"

func TestFormatDuplicatedFile(t *testing.T) {
	cases := []struct {
		p1, p2, want string
	}{
		{"/foo/bar/hello.x", "/foo/bar/world.y", "/foo/bar/[ hello.x == world.y ]"},
		{"/foo/bar/hello.x", "/foo/bar/zoo/world.y", "/foo/bar/[ hello.x == zoo/world.y ]"},
		{"/foo/bar/zoo/hello.x", "/foo/bar/world.y", "/foo/bar/[ zoo/hello.x == world.y ]"},
		{"/foo/bar/hello1", "/foo/bar/hello2/hello2", "/foo/bar/[ hello1 == hello2/hello2 ]"},
		{"/foo/bar/hello1/hello1", "/foo/bar/hello2", "/foo/bar/[ hello1/hello1 == hello2 ]"},
		{"/foo", "/bar", "/[ foo == bar ]"},
	}
	for _, c := range cases {
		got := FormatDuplicatedFile(c.p1, c.p2)
		if got != c.want {
			t.Errorf("FormatDuplicatedFile(%q, %q) == %q, want %q", c.p1, c.p2, got, c.want)
		}
	}
}
