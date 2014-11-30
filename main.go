// 2014 - Mathieu Lonjaret

// The acmetags program prints the tags of the acme windows.
package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"code.google.com/p/goplan9/plan9/acme"
)

var (
	allTags = flag.Bool("all", false, "print tags of all windows, instead of only \"win\" windows.")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: acmetags [-all]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func hostName() (string, error) {
	hostname, err := os.Hostname()
	if err == nil && hostname != "" {
		return hostname, nil
	}
	hostname = os.Getenv("HOSTNAME")
	if hostname != "" {
		return hostname, nil
	}
	out, err := exec.Command("hostname").Output()
	if err == nil && string(out) != "" {
		return strings.TrimSpace(string(out)), nil
	}
	return "", errors.New("all methods to find our hostname failed")
}

func main() {
	flag.Usage = usage
	flag.Parse()

	var hostname string
	var err error
	if !*allTags {
		hostname, err = hostName()
		if err != nil {
			log.Fatal(err)
		}
	}
	windows, err := acme.Windows()
	if err != nil {
		log.Fatal(err)
	}
	isWinHint := "-" + hostname
	for _, win := range windows {
		if !(*allTags || strings.HasSuffix(win.Name, isWinHint)) {
			continue
		}
		w, err := acme.Open(win.ID, nil)
		if err != nil {
			log.Fatalf("could not open window (%v, %d): %v", win.Name, win.ID, err)
		}
		tag, err := w.ReadAll("tag")
		if err != nil {
			log.Fatalf("could not read tags of window (%v, %d): %v", win.Name, win.ID, err)
		}
		fmt.Printf("%s\n\n", tag)
	}
}
