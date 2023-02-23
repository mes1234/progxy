package worker_test

import (
	"context"
	"testing"
	"time"

	"github.com/mes1234/progxy/internal/worker"
)

func TestShufflerShouldZeroAllBytes(t *testing.T) {

	//Arrange
	data := []byte{0xFF}
	sut := worker.NewShuffler(Generator(data), context.Background())

	//Act

	sut.Attach(ZeroAllBytes)

	time.Sleep(100 * time.Millisecond)
	//Assert
	for i, v := range data {
		if v != 0x00 {
			t.Fatalf("Expected all data to be 0x00 but got %d at position %d", v, i)
		}
	}

}

func TestShufflerShouldZeroAllAndAddOneBytes(t *testing.T) {

	//Arrange
	data := []byte{0xFF}
	sut := worker.NewShuffler(Generator(data), context.Background())

	//Act

	sut.Attach(ZeroAllBytes)
	sut.Attach(AddOne)

	//Assert
	time.Sleep(100 * time.Millisecond)

	for i, v := range data {
		if v != 0x01 {
			t.Fatalf("Expected all data to be 0x01 but got %d at position %d", v, i)
		}
	}
}

func TestShufflerShouldAndAddOneBytesInOrderedWay(t *testing.T) {

	//Arrange
	data := [][]byte{
		{0xAA, 0xAA},
		{0xBB, 0xBB, 0xBB},
		{0xCC, 0xCC, 0xCC, 0xCC},
		{0xDD, 0xDD, 0xDD, 0xDD, 0xDD},
	}

	copyOfData := [][]byte{
		{0xAA, 0xAA},
		{0xBB, 0xBB, 0xBB},
		{0xCC, 0xCC, 0xCC, 0xCC},
		{0xDD, 0xDD, 0xDD, 0xDD, 0xDD},
	}

	sut := worker.NewShuffler(IterativeGenerator(data), context.Background())

	//Act

	sut.Attach(AddOne)

	//Assert
	time.Sleep(100 * time.Millisecond)

	for j, row := range data {
		for i, v := range row {
			if v != copyOfData[j][i]+1 {
				t.Fatalf("Expected all data to be 0x01 but got %d at position %d", v, i)
			}
		}
	}

}

func Generator(data []byte) <-chan []byte {
	channel := make(chan []byte, 1)

	go func() {
		time.Sleep(1 * time.Millisecond)
		channel <- data
	}()

	return channel
}

func IterativeGenerator(data [][]byte) <-chan []byte {
	channel := make(chan []byte, 1)

	go func() {
		for _, v := range data {
			time.Sleep(1 * time.Millisecond)
			channel <- v
		}

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
