package worker

import (
	"context"
	"time"
)

type ProcessorFunc func(buffer []byte)

func CreateWriteToChannelProcessorFunc(channel chan<- []byte) ProcessorFunc {
	return func(buffer []byte) {
		channel <- buffer
	}
}

type Shufller interface {
	Attach(processor ProcessorFunc) error
}

type shuffler struct {
	processors    []ProcessorFunc
	myContext     context.Context
	parentContext context.Context
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

func (s *shuffler) start(input <-chan []byte) {
	select {
	case data := <-input:
		go s.processChunk(data)
	case <-s.myContext.Done():
		return
	}
}

func NewShuffler(input <-chan []byte, ctx context.Context) (Shufller, context.CancelFunc) {

	myContext, myCancelFunc := context.WithCancel(ctx)

	s := shuffler{
		processors:    make([]ProcessorFunc, 0),
		parentContext: ctx,
		myContext:     myContext,
	}

	go s.start(input)

	return &s, myCancelFunc
}
