package worker

import (
	"context"
	"time"
)

type ProcessorFunc func(buffer []byte)

type Shufller interface {
	Attach(processor ProcessorFunc) error
}

type shuffler struct {
	processors []ProcessorFunc
}

func (s *shuffler) Attach(processor ProcessorFunc) error {
	s.processors = append(s.processors, processor)
	return nil
}

// processChunk pass chunk of data to all processors one after another
// in ordered way
func (s *shuffler) processChunk(data []byte) {

	// wait  while processing might start before attachment of first Processor
	for len(s.processors) == 0 {
		time.Sleep(1 * time.Millisecond)
	}
	for _, p := range s.processors {
		p(data)
	}
}

func (s *shuffler) start(input <-chan []byte, ctx context.Context) {
	for {
		select {
		case data := <-input:
			s.processChunk(data)
		case <-ctx.Done():
			return
		}
	}

}

func NewShuffler(input <-chan []byte, ctx context.Context) Shufller {

	s := shuffler{
		processors: make([]ProcessorFunc, 0),
	}

	go s.start(input, ctx)

	return &s
}
