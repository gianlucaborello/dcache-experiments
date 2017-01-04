package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/etsy/statsd/examples/go"
)

type stat struct {
	entries int
	size    int64
}

func (s *stat) String() string {
	return fmt.Sprintf("entries: %d, size: %dMB",
		s.entries,
		s.size/1024/1024)
}

var stats stat

func walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	stats.entries++
	stats.size += info.Size()

	return nil
}

func main() {
	rootdir := flag.String("rootdir", "", "Root directory")

	flag.Parse()

	if *rootdir == "" {
		flag.Usage()
		os.Exit(1)
	}

	client := statsd.New("127.0.0.1", 8125)

	for {
		stats = stat{}
		t := time.Now()
		err := filepath.Walk(*rootdir, walkFn)
		if err != nil {
			log.Fatal(err)
		}

		diff := time.Since(t)

		client.Timing("worker.duration",
			diff.Nanoseconds()/int64(time.Millisecond))

		log.Print(stats.String())
		log.Printf("Elapsed: %v\n", diff)

		time.Sleep(5 * time.Second)
	}
}
