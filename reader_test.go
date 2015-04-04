package statos

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"testing"
)

func currentFile() string {
	_, filename, _, _ := runtime.Caller(1)
	return filename
}

func consumer(rs *ReaderStatos) chan bool {
	done := make(chan bool)

	go func() {
		ticker := time.Tick(1e9)
		for {
			bk := make([]byte, 64)
			// Throttle
			n, err := rs.Read(bk)
			<-ticker
			if n < 1 || err != nil {
				break
			}
		}
		done <- true
	}()

	return done
}

func progresser(rs *ReaderStatos, end chan bool) chan bool {
	done := make(chan bool)

	go func() {
		for {

			n, atEnd := rs.Progress()
			if atEnd {
				break
			}

			fmt.Printf("%v\r", n)

			select {
			case <-end:
				break
			default:
				continue
			}
		}
		done <- true
	}()

	return done
}

func TestReader(t *testing.T) {
	curFile := currentFile()
	r, err := os.Open(curFile)
	if err != nil {
		fmt.Printf("%s: %v\n", curFile, err)
		return
	}
	rs := NewReader(r)

	consumerChan := consumer(rs)
	done := progresser(rs, consumerChan)

	<-done

	defer rs.Close()
}
