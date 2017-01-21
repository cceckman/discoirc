// 2017-01-16 cceckman <charles@cceckman.com>
package socketmgr

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
	"path"
	"strings"
	"testing"
)

const (
	evname = "TEST_ENVVAR"
)

func listenCases(tmpdir string) []*sm {
	return []*sm{
		&sm{
			target: "",
			paths:  []string{path.Join(tmpdir, "subdir", "asocket")},
		},
		&sm{
			target: "",
			paths: []string{
				path.Join(tmpdir, fmt.Sprintf("$%s", evname), "socket"),
			},
		},
		&sm{
			target: path.Join(tmpdir, "subdir", "bsocket"),
			paths: []string{
				path.Join(tmpdir, fmt.Sprintf("$%s", evname), "socket"),
			},
		},
	}
}

// Test that the permissions on the socket are user-only, or that the socket is nil.
func testPerms(t *testing.T, pth string) {
	if pth == "" {
		return
	}

	info, err := os.Stat(pth)
	if err != nil {
		t.Errorf("could not stat listener path %s: %v", pth, err)
		return
	}

	var userbits os.FileMode = 0750
	perms := info.Mode().Perm()
	if perms != userbits {
		t.Errorf("file %s has wrong mode: got: %o want: %o", info.Name(), perms, userbits)
	}
}

// setupListen sets up a test environment for Listen tests and returns a temp directory.
// When the 'done' channel is closed, it cleans up the test environment.
func setupListen(t *testing.T, done chan struct{}) (string, error) {
	tmpdir, err := ioutil.TempDir("", "discoirc-testing")
	if err != nil {
		return "", fmt.Errorf("error in creating empty directory: %v", err)
	}
	go func() {
		<-done
		if err := os.RemoveAll(tmpdir); err != nil {
			t.Fatalf("could not clean up %s: %v", tmpdir, err)
		}
	}()

	if err := os.Setenv(evname, evname); err != nil {
		return "", fmt.Errorf("error in creating empty directory: %v", err)
	}
	go func() {
		<-done
		if err := os.Unsetenv(evname); err != nil {
			t.Fatalf("could not set environment variable %s: %v", evname, err)
		}
	}()

	return tmpdir, nil
}

// Test that Listen works (fails) when paths do not exist.
func TestListenPathDne(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	tmpdir, err := setupListen(t, done)
	if err != nil {
		t.Fatal(err)
	}

	// Directories / paths that don't exist.
	listenCases := listenCases(tmpdir)
	// Cases that we expect to return an error.
	for _, s := range listenCases {
		lis, path, err := s.Listen()
		testPerms(t, path)
		if err == nil {
			t.Errorf("expected no resolution for %v, got %s", s.paths, path)
			lis.Close()
		}
	}
}

// Test that Listen works (succeeds) on paths that exist exactly.
func TestListenExplicit(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	tmpdir, err := setupListen(t, done)
	if err != nil {
		t.Fatal(err)
	}

	listenCases := listenCases(tmpdir)

	// Create directories.
	for _, s := range listenCases {
		// Create existing directories.
		d := []string{path.Dir(s.target)}
		for _, p := range s.paths {
			d = append(d, path.Dir(p))
		}
		SetupDirs(d)
	}

	// Create an explicitly-named socket in the existing directory.
	for _, s := range listenCases {
		want := s.target
		if want == "" {
			want = s.paths[0]
		}
		want = os.ExpandEnv(want)

		lis, got, err := s.Listen()
		testPerms(t, got)
		if err != nil {
			t.Errorf("expected resolution for %v, got error %v", s.paths, err)
		}
		if got != want {
			t.Errorf("unexpected socket path: got: %s want: %s", got, want)
		}
		if lis == nil {
			t.Errorf("unexpected nil listener: lis %v path %s err %v", lis, got, err)
		} else {
			lis.Close()
		}
	}
}

// Test that Listen works (succeeds) on paths that are existing directories.
func TestListenDirectory(t *testing.T) {
	done := make(chan struct{})
	defer close(done)
	tmpdir, err := setupListen(t, done)
	if err != nil {
		t.Fatal(err)
	}

	listenCases := listenCases(tmpdir)

	// Create directories.
	for _, s := range listenCases {
		// Create existing directories.
		d := []string{path.Dir(s.target)}
		for _, p := range s.paths {
			d = append(d, path.Dir(p))
		}
		SetupDirs(d)
	}

	// Substitute directories for paths.
	for _, s := range listenCases {
		for i, p := range s.paths {
			s.paths[i] = path.Dir(p)
		}
	}

	// Cases where we create a random socket in an existing directory.
	for _, s := range listenCases {
		wantPrefix := path.Dir(s.target)
		if s.target == "" {
			wantPrefix = s.paths[0]
		}
		wantPrefix = os.ExpandEnv(wantPrefix)

		lis, got, err := s.Listen()
		testPerms(t, got)
		if err != nil {
			t.Errorf("expected resolution for %v, got error %v", s.paths, err)
		}
		if !strings.HasPrefix(got, wantPrefix) {
			t.Errorf("unexpected socket path: got: %s want prefix: %s", got, wantPrefix)
		}
		if lis == nil {
			t.Errorf("unexpected nil listener: lis %v path %s err %v", lis, got, err)
		} else {
			lis.Close()
		}
	}
}

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
		// No target; empty directory.
		&sm{
			target: "",
			paths:  []string{"$ENVVARDNE"},
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
	// Don't check permissions, created just for testing.
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
					"$ENVVARDNE",
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
