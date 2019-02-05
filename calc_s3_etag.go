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
	// flags
	chunkSizeMbPtr := flag.Int("chunksize", 8, "Multi-part chunksize in MB used for upload")
	includeFNPtr := flag.Bool("fn", false, "include filename in output")
	includeFNAPtr := flag.Bool("fna", false, "include filename in output after etag")
	md5ForSingleMultipart := flag.Bool("md5_for_single_multipart", false, "just print MD5 sum for files under chunksize, s3cmd upload functionality")
	s3cmdStylePtr := flag.Bool("s3cmd_style", false, "s3cmd style: set both md5_for_single_multipart and chunksize=15")
	flag.Parse()
	args := flag.Args()
	filename := args[0]
	if *s3cmdStylePtr == true {
		*chunkSizeMbPtr = 15
		*md5ForSingleMultipart = true
	}

	isThere, _ := exists(filename)
	if isThere {
		etag := GetEtag(filename, *chunkSizeMbPtr, *md5ForSingleMultipart)

		if *includeFNPtr == true {
			fmt.Print(filename, ": ")
		}

		fmt.Print(etag)

		if *includeFNAPtr == true {
			fmt.Print("  ", filename)
		}
		fmt.Println()
	} else {
		fmt.Println("Could not find file \"", filename, "\"")
	}
}

func GetEtag(path string, chunkSizeMb int, md5ForSingleMultipart bool) string {
	etag := ""
	chunkSize := chunkSizeMb * 1024 * 1024

	f, err := os.Open(path)
	check(err)
	defer f.Close()

	r := bufio.NewReader(f)

	parts := 0
	chunk := make([]byte, chunkSize)
	md5list := make([]byte, 0)
	for {
		bytesRead, err := r.Read(chunk)
		if err == io.EOF {
			break
		}
		check(err)

		parts += 1
		chunkHash := md5.Sum(chunk[0:bytesRead])
		//fmt.Println("bytesRead: ", bytesRead, " chunkHash: ", chunkHash)
		md5list = append(md5list, chunkHash[:]...)
	}

	if parts < 2 && md5ForSingleMultipart == true {
		if parts == 1 {
			etag = fmt.Sprintf("%x", md5list)
		} else { // 0 size / 0 parts
			etag = fmt.Sprintf("%x", md5.Sum(chunk[0:0]))
		}
	} else {
		listHash := md5.Sum(md5list)
		etag = fmt.Sprintf("%x-%d", listHash, parts)
	}

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
