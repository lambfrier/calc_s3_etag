// Standalone utility to calculate AWS multi-part uploaded S3 file Etags
//
// Inspired from Antonio Espinosa & r03's SO answer to:
// https://stackoverflow.com/questions/6591047/etag-definition-changed-in-amazon-s3

package main

import (
	"bufio"
	"crypto/md5"
	"flag"
	"fmt"
	"io"
	"os"
)

// Main
func main() {
	chunkSizeMbPtr := flag.Int("chunksize", 8, "Multi-part chunksize in MB used for upload")
	includeFNPtr := flag.Bool("fn", false, "include filename in output")
	flag.Parse()
	args := flag.Args()
	filename := args[0]

	isThere, _ := exists(filename)
	if isThere {
		etag := GetEtag(filename, *chunkSizeMbPtr)
		if *includeFNPtr == true {
			fmt.Print(filename + ": ")
		}
		fmt.Println(etag)
	} else {
		fmt.Println("Could not find file \"" + filename + "\"")
	}
}

func GetEtag(path string, chunkSizeMb int) string {
	chunkSize := chunkSizeMb * 1024 * 1024

	f, err := os.Open(path)
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	parts := 0
	chunk := make([]byte, chunkSize)
	contentToHash := make([]byte, 0)
	for {
		bytesRead, err := r.Read(chunk)
		if err == io.EOF {
			break
		}
		check(err)

		parts += 1
		hash := md5.Sum(chunk[0:bytesRead])
		//fmt.Println("bytesRead: ", bytesRead, " hash: ", hash)
		contentToHash = append(contentToHash, hash[:]...)
	}

	hash := md5.Sum(contentToHash)
	etag := fmt.Sprintf("%x-%d", hash, parts)

	return etag
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func check(e error) {
	if e != nil {
		panic(fmt.Sprintf("Error encountered - ", e))
	}
}
