package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/radovskyb/watcher"
)

func main() {
	w1 := watcher.New()
	w2 := watcher.New()

	testdir1 := "/mnt/ssd2/testfile"
	testdir2 := "/mnt/ssd2/testfile2"

	go func() {
		for {
			select {
			case ev := <-w1.Event:
				if ev.Op == watcher.Move {
					// fmt.Println("found watch1 move: " + ev.Path[66:])
					// fmt.Println("found watch1 move old: " + ev.OldPath[66:])
					_ = os.Rename(
						testdir2+ev.OldPath[len(testdir1):],
						testdir2+ev.Path[len(testdir1):],
					)
				} else if ev.Op == watcher.Remove {
					// fmt.Println("found watch1 remove: " + ev.Path)
				}
			case err := <-w1.Error:
				log.Fatalln(err)
			case <-w1.Closed:
				return
			}
		}
	}()

	go func() {
		for {
			select {
			case ev := <-w2.Event:
				if ev.Op == watcher.Move {
					// fmt.Println("found watch2 move: " + ev.Path)
					// fmt.Println("found watch2 move old: " + ev.OldPath)
				} else if ev.Op == watcher.Remove {
					fmt.Println("found watch2 remove: " + ev.Path)
				}
			case err := <-w2.Error:
				log.Fatalln(err)
			case <-w2.Closed:
				return
			}
		}
	}()

	// // Watch this folder for changes.
	// if err := w1.Add("."); err != nil {
	// 	log.Fatalln(err)
	// }
	// if err := w2.Add("."); err != nil {
	// 	log.Fatalln(err)
	// }

	// Watch test_folder recursively for changes.
	if err := w1.AddRecursive(testdir1); err != nil {
		log.Fatalln(err)
	}
	if err := w2.AddRecursive(testdir2); err != nil {
		log.Fatalln(err)
	}

	// Print a list of all of the files and folders currently
	// being watched and their paths.
	// for path, f := range w2.WatchedFiles() {
	// 	fmt.Printf("%s: %s\n", path, f.Name())
	// }

	// fmt.Println()

	// // Trigger 2 events after watcher started.
	// go func() {
	// 	w1.Wait()
	// 	w1.TriggerEvent(watcher.Create, nil)
	// 	w1.TriggerEvent(watcher.Remove, nil)
	// }()

	// Start the watching process - it'll check for changes every 100ms.
	fmt.Println("w1 starting...")
	go func() {
		if err := w1.Start(time.Millisecond * 100); err != nil {
			log.Fatalln(err)
		}
	}()
	fmt.Println("w1 start")
	fmt.Println("w2 starting...")
	if err := w2.Start(time.Millisecond * 100); err != nil {
		log.Fatalln(err)
	}
	fmt.Println("w2 start")
}
