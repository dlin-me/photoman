package main

import (
	"testing"
	"os"
	"time"
)

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

func TestGetExifDateTime(t *testing.T) {
	want := "Sat Sep  2 14:30:10 AEDT 2000"
	pwd, _ := os.Getwd();
	tm, _ := GetExifDateTime(pwd+"/test/test.jpg");
	got := tm.Format(time.UnixDate)
	if got != want {
		t.Errorf("GetExifDateTime() = %q, want %q", got, want)
	}
}

func TestGetProposedPath(t  *testing.T) {
	tm, _ := time.Parse("20060102150405", "20150102030405")
	path, _ := os.Getwd()
	want := path + "/2015/2015_01"
	got := GetProposedPath(tm)
	if got != want {
		t.Errorf("GetProposedPath() = %q, want %q", got, want)
	}
}
