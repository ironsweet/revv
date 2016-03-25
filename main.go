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

	//cmd := exec.Command(os.Args[1], os.Args[1:]...)
	//cmd.Stdout = os.Stdout
	//cmd.Stderr = os.Stderr
	//err := cmd.Start()
	//if err != nil {
	//	log.Fatal(err)
	//}

	done := make(chan bool)
	go func() {
		var timer *time.Timer
		for {
			select {
			case event := <-watcher.Events:
				log.Println("event:", event)
				if event.Op&fsnotify.Write == fsnotify.Write {
					if timer == nil {
						timer = time.NewTimer(3 * time.Second)
						go func() {
							<-timer.C
							fmt.Println("restart")
						}()
					} else {
						timer.Reset(5 * time.Second)
					}
				}
			case err := <-watcher.Errors:
				log.Println("error:", err)
			}
		}
	}()

	p, err := exec.LookPath(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Monitoring: %v", p)

	err = watcher.Add(p)
	if err != nil {
		log.Fatal(err)
	}
	<-done
}
