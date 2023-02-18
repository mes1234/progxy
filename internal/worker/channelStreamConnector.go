package worker

import (
	"fmt"
	"io"
	"sync"
)

// Should be used as goroutine otherwise it will never release thread
func readAndForward(out channelWrapper, reader io.Reader, wg *sync.WaitGroup) {
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
			out.channel <- outBuf
		}

	}
}

// Should be used as goroutine otherwise it will never release thread
func forwardToWriter(in channelWrapper, writer io.Writer) {
	for {
		data := <-in.channel
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
func CreateChannelFromReaderWriter(description string, rw io.ReadWriter) (in chan []byte, out chan []byte, wg *sync.WaitGroup) {

	wg = &sync.WaitGroup{}
	wg.Add(1)

	outWrapper := channelWrapper{
		channel:     make(chan []byte, 1),
		description: "out channel for " + description,
	}
	inWrapper := channelWrapper{
		channel:     make(chan []byte, 1),
		description: "in channel for " + description,
	}

	// Read data from Reader and push to channel
	go readAndForward(outWrapper, rw, wg)

	// Forward data from in channel to Writer
	go forwardToWriter(inWrapper, rw)

	return inWrapper.channel, outWrapper.channel, wg
}

type channelWrapper struct {
	channel     chan []byte
	description string
}
