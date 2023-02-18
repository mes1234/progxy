package worker

import (
	"fmt"
	"io"
	"sync"
)

// Should be used as goroutine otherwise it will never release thread
func readAndForward(out chan<- []byte, reader io.Reader, wg *sync.WaitGroup) {
	defer wg.Done()
	readBuf := make([]byte, bufferSize)
	for {
		n, err := reader.Read(readBuf)
		if err != nil {
			return
		}
		if n != 0 {
			outBuf := make([]byte, n)
			copy(outBuf, readBuf)
			out <- outBuf
		}

	}
}

// Should be used as goroutine otherwise it will never release thread
func forwardToWriter(in <-chan []byte, writer io.Writer) {
	for {
		data := <-in
		n, err := writer.Write(data)
		if err != nil {
			fmt.Printf("read %v data and failed", n)
			return
		}
	}
}

// CreateChannelFromReaderWriter
// performs Read on rw and push it to out chan
// retrieve from in chan and performs Write
func CreateChannelFromReaderWriter(rw io.ReadWriter) (in chan []byte, out chan []byte, wg *sync.WaitGroup) {

	wg = &sync.WaitGroup{}
	wg.Add(1)

	out = make(chan []byte, 1)
	in = make(chan []byte, 1)

	// Read data from Reader and push to channel
	go readAndForward(out, rw, wg)

	// Forward data from in channel to Writer
	go forwardToWriter(in, rw)

	return in, out, wg
}
