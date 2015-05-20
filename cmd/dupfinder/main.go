package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) <= 1 {
		log.Fatalln("Usage: dupfinder <directory>")
	}

	directoryPath := os.Args[1]

	fi, err := os.Stat(directoryPath)
	if err != nil {
		log.Fatalln("Error:", err)
	}

	if !fi.IsDir() {
		log.Fatalf("Not a directory: %q", fi.Name())
	}

	filemap := map[string][]string{}

	if err := filepath.Walk(directoryPath, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			if info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		f, err := os.Open(path)
		if err != nil {
			log.Printf("Error opening file: %q: %v", info.Name(), err)
			return nil
		}

		h := md5.New()
		if _, err := io.Copy(h, f); err != nil {
			log.Printf("Error reading file: %q", info.Name())
			return nil
		}

		hashValue := fmt.Sprintf("%x", h.Sum(nil))

		strArr, ok := filemap[hashValue]
		if !ok {
			strArr = []string{}
		}
		strArr = append(strArr, f.Name())
		filemap[hashValue] = strArr

		return nil
	}); err != nil {
		log.Fatalln("Error:", err)
	}

	dupsFound := false
	counter := 1
	for _, strs := range filemap {
		if len(strs) > 1 {
			if !dupsFound {
				log.Println("Duplicates found!")
				dupsFound = true
			}

			log.Printf("%d. Duplicate files: [%s]",
				counter,
				strings.Join(strs, ", "))
			counter++
		}
	}

	if !dupsFound {
		log.Println("No duplicates found.")
	}
}
