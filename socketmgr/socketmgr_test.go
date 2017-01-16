// 2017-01-16 cceckman <charles@cceckman.com>
package socketmgr

import (
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"testing"
)

func TestGet(t *testing.T) {
	tmpdir, err := ioutil.TempDir("", "discoirc-testing")
	if err != nil {
		t.Fatal("error in creating empty directory: ", err)
	}
	defer func() {
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Fatalf("could not clean up %s: %v", tmpdir, err)
		}
	}()

	// Cases that we expect to resolve nothing.
	for _, s := range []*sm{
		// No target; directory doesn't exist.
		&sm{
			target: "",
			paths:  []string{path.Join(tmpdir, "does-not-exist")},
		},
		// No target; directory has no sockets.
		&sm{
			target: "",
			paths:  []string{tmpdir},
		},
	} {
		got := s.Get()
		if len(got) > 0 {
			t.Errorf("expected no resolution for %v, got %s", s, strings.Join(got, ","))
		}
	}

	// Create a socket.
	sockaddr := path.Join(tmpdir, "test.sock")
	lis, err := net.Listen("unix", sockaddr)
	if err != nil {
		t.Fatalf("failed to listen: %v", err)
	}
	defer lis.Close()

	for _, cs := range []struct {
		Mgr  T
		Want []string
	}{
		{
			Mgr: &sm{
				target: sockaddr,
				paths:  []string{},
			},
			Want: []string{sockaddr},
		},
		{
			Mgr: &sm{
				target: "",
				paths: []string{
					path.Join(tmpdir, "does-no-exist"),
					tmpdir,
				},
			},
			Want: []string{sockaddr},
		},
	} {
		got := cs.Mgr.Get()
		want := cs.Want
		if len(got) != len(want) {
			t.Errorf("got: %v want: %v", got, want)
		}
		for i := range got {
			if want[i] != got[i] {
				t.Errorf("got: %v want: %v", got, want)
			}
		}
	}
}
