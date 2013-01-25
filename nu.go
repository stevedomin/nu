package main

// Naïve approach

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime/pprof"
	"strings"
)

var (
	cpuprofile    = flag.String("cpuprofile", "", "write cpu profile to file")
	memprofile    = flag.String("memprofile", "", "write memory profile to this file")
	patternGlobal string
)

func main() {
	flag.Parse()

	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.WriteHeapProfile(f)
		f.Close()
		return
	}

	args := flag.Args()
	if len(args) != 2 {
		fmt.Println("Usage: nu <pattern> <directory>")
		os.Exit(1)
	}

	fmt.Printf("Search : %s \n\n", args[0])
	nu(args[0], args[1])
}

func nu(pattern, directory string) {
	patternGlobal = pattern

	filepath.Walk(directory, walkFn)
}

func walkFn(path string, info os.FileInfo, err error) error {

	if info.IsDir() == true {
		return nil
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Can't open %s", f)
		return nil
	}
	defer f.Close()

	data := make([]byte, 256)
	count, err := f.Read(data)
	isBinary := bytes.Contains(data[:count], []byte{0})
	f.Seek(0, 0)
	if isBinary {
		return nil
	}

	bf := bufio.NewReader(f)
	var lineNumber int64 = 0
	fileMatch := false
	headerPrinted := false

	for {
		lineNumber++

		line, isPrefix, err := bf.ReadLine()
		if err == io.EOF {
			if fileMatch {
				fmt.Println("")
			}
			return nil
		}

		if err == nil && !isPrefix {
			if strings.Contains(string(line), patternGlobal) {
				fileMatch = true

				if fileMatch && !headerPrinted {
					fmt.Printf("\x1b[1m%s\x1b[0m\n", path)
					headerPrinted = true
				}

				fmt.Printf("\x1b[1m%b\x1b[0m: 	%s \n", lineNumber, line)
			}
		}
	}

	return nil
}
