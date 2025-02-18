package views

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"testing"
	"time"
)

func TestWatchDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "test_watch_dir*")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan bool)
	events, err := watchDir(dir, regexp.MustCompile("(?i)^[a-z0-9].+html"), done)
	if err != nil {
		t.Fatal(err)
	}

	var counter int = 0
	go func() {
		for range events {
			counter++
		}
	}()

	for _, newfile := range []string{
		"_newfile.html", // not counted
		"abc.html",      // counted
		"abc.HTML",      // counted
		".hidden.HTML",  // not counted
	} {
		o, err := os.Create(filepath.Join(dir, newfile))
		if err != nil {
			t.Fatal(err)
		}
		_, err = fmt.Fprint(o, "hi")
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(50 * time.Millisecond)
	done <- true

	_ := os.RemoveAll(dir)

	if got, want := counter, 2; got != want {
		t.Errorf("counter got %d want %d", got, want)
	}
}

func TestWatchDirFunc(t *testing.T) {

	dir, err := os.MkdirTemp("", "test_watch_dir_func*")
	if err != nil {
		t.Fatal(err)
	}

	done := make(chan bool)
	var counter int = 0
	counterFunc := func() error {
		counter++
		return nil
	}
	go watchDirFunc(dir, regexp.MustCompile("(?i)^[a-z0-9].+html"), done, counterFunc)
	if err != nil {
		t.Fatal(err)
	}

	for _, newfile := range []string{
		"_newfile.html", // not counted
		"abc.html",      // counted
		"abc.HTML",      // counted
		".hidden.HTML",  // not counted
	} {
		o, err := os.Create(filepath.Join(dir, newfile))
		if err != nil {
			t.Fatal(err)
		}
		_, err = fmt.Fprint(o, "hi")
		if err != nil {
			t.Fatal(err)
		}
	}

	time.Sleep(50 * time.Millisecond)
	done <- true

	_ := os.RemoveAll(dir)

	if got, want := counter, 2; got != want {
		t.Errorf("counter got %d want %d", got, want)
	}
}
