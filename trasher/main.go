package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/satori/go.uuid"
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

	numfiles := 0

	for {
		uuid := uuid.NewV4()
		os.Open(path.Join(*rootdir, uuid.String()))
		numfiles++
		if numfiles%1000000 == 0 {
			log.Printf("Number of processed files: %d\n", numfiles)
		}
	}
}
