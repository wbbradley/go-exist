// A thread-safe mechanism for swapping out existence keys.

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

// Read keys from a stream delimited by line breaks
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

// Read keys from a slice
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
