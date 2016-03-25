package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"os/exec"

	"github.com/fsnotify/fsnotify"
)

func main() {
	if len(os.Args) == 1 {
		fmt.Println("revv auto restarts program when it detects binary change.")
		fmt.Println("Usage: revv <program> [args...]")
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	cmd := exec.Command(os.Args[1], os.Args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	done := make(chan bool)
	go func() {
		var timer *time.Timer
		for {
			select {
			case event := <-watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Println("Detected change.")
					if timer == nil {
						timer = time.NewTimer(3 * time.Second)
						go func() {
							<-timer.C
							log.Println("Reloading...")
							if err = cmd.Process.Kill(); err != nil {
								log.Fatal("Kill:", err)
							} else {
								cmd.Wait() // wait until program ends
								cmd = exec.Command(os.Args[1], os.Args[1:]...)
								cmd.Stdout = os.Stdout
								cmd.Stderr = os.Stderr
								if err = cmd.Start(); err != nil {
									log.Fatal("Start", err)
								}
								log.Println("Reloaded.")
							}
							timer = nil
						}()
					} else {
						timer.Reset(5 * time.Second)
					}
				}
			case err := <-watcher.Errors:
				log.Println("Error:", err)
			}
		}
	}()

	var p string
	if p, err = exec.LookPath(os.Args[1]); err != nil {
		log.Fatal(err)
	}
	log.Println("Monitoring:", p)
	if err = watcher.Add(p); err != nil {
		log.Fatal(err)
	}
	<-done
}
