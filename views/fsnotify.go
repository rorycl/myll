package views

import (
	"fmt"
	"path/filepath"
	"regexp"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// flushDuration sets the time given to wait for multiple editor writes
var flushDuration time.Duration = 50 * time.Millisecond

// watchDir watches a directory for events, deferring the monitor the
// fileLoop goroutine. This code is largely taken from:
// https://github.com/fsnotify/fsnotify/blob/v1.8.0/cmd/fsnotify/file.go
func watchDir(path string, fileMatcher *regexp.Regexp, done <-chan bool) (<-chan bool, error) {

	eventChan := make(chan bool)
	bufferChan := make(chan bool)

	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("fsnotify create error: %v", err)
	}

	err = w.Add(path)
	if err != nil {
		return nil, fmt.Errorf("fsnotify add error: %v", err)
	}

	// Start listening for events.
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-done:
				return
			case err, ok := <-w.Errors:
				if !ok {
					panic("unexpected close")
				}
				panic(fmt.Sprintf("unexpected notify error: %v", err))

			case e, ok := <-w.Events:
				if !ok { // Watcher.Close() called?
					panic("unexpected close")
					return
				}
				if !e.Has(fsnotify.Write) {
					continue
				}
				// fmt.Println(e.Name, e.String())
				basename := filepath.Base(e.Name)
				if fileMatcher.MatchString(basename) {
					eventChan <- true
				}
			}
		}
	}()

	// crude buffer of double writes by editors like vim
	go func() {
		flush := false
		timer := time.NewTicker(flushDuration)
		for {
			select {
			case <-eventChan:
				flush = true
			case <-timer.C:
				if flush {
					bufferChan <- true
					flush = false
				}
			}
		}
	}()

	go func() {
		wg.Wait()
		close(eventChan)
		defer w.Close()
	}()

	return bufferChan, nil
}

/*
func main() {
	done := make(chan bool)
	events, err := watchDir("/tmp/zan", regexp.MustCompile("(?i)^[a-z0-9].+html$"), done)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _ = range events {
		fmt.Println("got!")
	}
}
*/
