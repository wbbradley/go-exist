// The goal of this library is to enable you to create a thread-safe mechanism
// for answering questions of string membership in sets.  Altering what
// "exists" can be done by appending at run-time, or by transacting a swap of a
// wholly new set.
package filter

import (
	"bufio"
	"io"
	"log"
	"os"
)

type ExistenceFilter interface {
	KeyExists(key string) bool
}

type MutableExistenceFilter interface {
	ExistenceFilter
	ImportKeys(keyCount uint, keys <-chan string)
}

// Read keys from a stream delimited by line breaks. These keys are passed into
// a channel which is then consumed by ImportKeys.
func ReadStreamIntoFilter(f MutableExistenceFilter, keyCount uint, r io.Reader) {
	keys := make(chan string, 100)
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			if key := scanner.Text(); len(key) > 0 {
				keys <- key
			}
		}
		close(keys)
		if err := scanner.Err(); err != nil {
			log.Println("ReadStreamIntoFilter - scanner error", err)
			if wd, err := os.Getwd(); err == nil {
				log.Println("Current working directory is", wd)
			}
			return
		}
	}()

	f.ImportKeys(keyCount, keys)
}

// Read keys from a slice into a MutableExistenceFilter. This function converts
// a slice of strings into a channel of strings.  It then simply passes that
// channel on to ImportKeys.
func ReadKeysIntoFilter(f MutableExistenceFilter, keysToRead []string) {
	keys := make(chan string, 100)
	go func() {
		for _, key := range keysToRead {
			if len(key) > 0 {
				keys <- key
			}
		}
		close(keys)
	}()

	f.ImportKeys(uint(len(keysToRead)), keys)
}
