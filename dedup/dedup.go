/*
Copyright 2015 Mathieu Lonjaret

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

     http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// The dedup program deduplicates to stdout the acmetags-dump files in the current directory.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var flagVerbose = flag.Bool("v", false, "be verbose")

func usage() {
	fmt.Fprintf(os.Stderr, "usage: dedup\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()

	dedup := make(map[string]struct{}) // unused in verbose mode

	// for verbose mode
	var (
		totalFiles   int
		totalLines   int
		duplicates   int
		distribution = make(map[string]int)
	)

	d, err := os.Open(".")
	if err != nil {
		log.Fatal(err)
	}
	defer d.Close()
	dir, err := d.Readdirnames(-1)
	if err != nil {
		log.Fatal(err)
	}
	for _, filename := range dir {
		if !strings.HasPrefix(filename, "acmetags.dump-") {
			continue
		}
		if *flagVerbose {
			totalFiles++
			log.Printf("processing %v", filename)
		}
		f, err := os.Open(filename)
		if err != nil {
			log.Fatal(err)
		}
		sc := bufio.NewScanner(f)
		for sc.Scan() {
			if *flagVerbose {
				totalLines++
				l := sc.Text()
				n, ok := distribution[l]
				if ok {
					duplicates++
					n++
					log.Printf("duplicate: %q", l)
				}
				distribution[l] = n
			} else {
				// TODO(mpl): is it faster to check and not overwrite if already exists ?
				dedup[sc.Text()] = struct{}{}
			}
		}
		if err := sc.Err(); err != nil {
			log.Fatal(err)
		}
		f.Close()
	}
	if *flagVerbose {
		for k, v := range distribution {
			// TODO(mpl): buffer ?
			fmt.Printf("[%s]: %d\n", k, v)
		}
	} else {
		for k, _ := range dedup {
			// TODO(mpl): buffer ?
			fmt.Println(k)
		}
	}
	if *flagVerbose {
		log.Printf("total files processed: %d", totalFiles)
		log.Printf("total lines: %d", totalLines)
		log.Printf("total duplicates: %d", duplicates)
	}
}
