// 2017-01-16 cceckman <charles@cceckman.com>
package socketmgr

import (
	"io/ioutil"
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
}
