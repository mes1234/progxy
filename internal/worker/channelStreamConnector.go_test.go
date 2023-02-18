package worker_test

import (
	"io"
	"testing"
	"time"

	"github.com/mes1234/progxy/internal/worker"
)

func TestStreamShouldReadOneTimeAndPutItToOutChanTime(t *testing.T) {
	//Arrange

	dataChunk := [][]byte{
		{0xAA, 0xAA},
	}

	dummyRW := newDummyStreamReadWriter(dataChunk)

	//Act
	_, outChan, _ := worker.CreateChannelFromReaderWriter(dummyRW)

	//Assert
	data := <-outChan
	for i, v := range data {
		if v != dataChunk[0][i] {
			t.Fatalf("Expected to retrieve %v at postion %v but got %v", dataChunk[0][i], i, v)
		}
	}
}

func TestStreamShouldReadMoreThanOneTimeAndPutItToOutChanTime(t *testing.T) {
	//Arrange

	dataChunk := [][]byte{
		{0xAA, 0xAA},
		{0xBB, 0xBB, 0xBB},
		{0xCC, 0xCC, 0xCC, 0xCC},
		{0xDD, 0xDD, 0xDD, 0xDD, 0xDD},
	}

	dummyRW := newDummyStreamReadWriter(dataChunk)

	//Act
	_, outChan, _ := worker.CreateChannelFromReaderWriter(dummyRW)

	time.Sleep(1 * time.Second)

	//Assert
	for j, dataChunkRow := range dataChunk {
		data := <-outChan
		for i, v := range data {
			if v != dataChunkRow[i] {
				t.Fatalf("Expected to retrieve %v at postion %v but got %v for data row %v", dataChunkRow[i], i, v, j)
			}
		}
	}
}

func TestStreamShouldGetDataFromInChannelAndPassItToWriter(t *testing.T) {
	//Arrange

	dataChunk := [][]byte{
		{0xAA, 0xAA},
		{0xBB, 0xBB, 0xBB},
		{0xCC, 0xCC, 0xCC, 0xCC},
		{0xDD, 0xDD, 0xDD, 0xDD, 0xDD},
	}
	//Act
	dummyRW := newDummyStreamReadWriterWithAssertions(dataChunk, t)
	inChan, _, _ := worker.CreateChannelFromReaderWriter(dummyRW)

	for _, dataChunkRow := range dataChunk {
		inChan <- dataChunkRow
	}

	//Assert
	//Done by Write metod directly

}

type dummyStreamReadWriter struct {
	readsCount      int
	writesCount     int
	dataChunk       [][]byte
	assertionHelper *testing.T
}

func newDummyStreamReadWriter(dataChunks [][]byte) io.ReadWriter {
	return &dummyStreamReadWriter{
		readsCount: len(dataChunks),
		dataChunk:  dataChunks,
	}
}

func newDummyStreamReadWriterWithAssertions(dataChunks [][]byte, t *testing.T) io.ReadWriter {
	return &dummyStreamReadWriter{
		writesCount:     len(dataChunks),
		dataChunk:       dataChunks,
		assertionHelper: t,
	}
}

func (dsr *dummyStreamReadWriter) Read(p []byte) (n int, err error) {
	if dsr.readsCount > 0 {
		data := dsr.dataChunk[len(dsr.dataChunk)-dsr.readsCount]
		copy(p, data)
		dsr.readsCount--
		return len(data), nil
	}
	return 0, io.EOF
}

func (dsr *dummyStreamReadWriter) Write(p []byte) (n int, err error) {
	if dsr.writesCount > 0 {
		data := dsr.dataChunk[len(dsr.dataChunk)-dsr.writesCount]
		time.Sleep(1 * time.Second)
		if len(data) != len(p) {
			dsr.assertionHelper.Fatalf("Expected to get array of len %v but got %v", len(data), len(p))
		}
		dsr.writesCount--
		return len(p), nil
	}

	return 0, io.EOF
}
