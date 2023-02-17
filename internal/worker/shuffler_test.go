package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/mes1234/progxy/internal/worker"
)

func TestShufflerShouldZeroAllBytes(t *testing.T) {

	//Arrange
	data := make([]byte, 1024)
	sut, _ := worker.NewShuffler(Generator(data, 0xFF), context.Background())

	//Act

	sut.Attach(ZeroAllBytes)

	time.Sleep(10 * time.Millisecond)
	//Assert
	for i, v := range data {
		if v != 0x00 {
			t.Fatalf("Expected all data to be 0x00 but got %d at position %d", v, i)
		}
	}

}

func TestShufflerShouldZeroAllAndAddOneBytes(t *testing.T) {

	//Arrange
	data := make([]byte, 1024)
	sut, _ := worker.NewShuffler(Generator(data, 0xFF), context.Background())

	//Act

	sut.Attach(ZeroAllBytes)
	sut.Attach(AddOne)

	//Assert
	time.Sleep(10 * time.Millisecond)

	for i, v := range data {
		if v != 0x01 {
			t.Fatalf("Expected all data to be 0x01 but got %d at position %d", v, i)
		}
	}

}

func Generator(data []byte, value byte) <-chan []byte {
	channel := make(chan []byte, 1)
	for i := range data {
		data[i] = value
	}
	go func() {
		time.Sleep(1 * time.Millisecond)
		channel <- data
	}()

	return channel
}

func ZeroAllBytes(buffer []byte) {
	for i := range buffer {
		buffer[i] = 0x00
	}
}

func AddOne(buffer []byte) {
	for i := range buffer {
		buffer[i] = buffer[i] + 0x01
	}
}
