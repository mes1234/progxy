package worker

import "context"

type ProcessorFunc func(buffer []byte)

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

func (s *shuffler) processChunk(data []byte) {
	for _, p := range s.processors {
		p(data)
	}
}

func (s *shuffler) start(input <-chan []byte) {
	select {
	case data := <-input:
		// trigger processing and carry on
		go s.processChunk(data)
	case <-s.myContext.Done():
		//Cancellation of shuffler
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

	s.start(input)

	return &s, myCancelFunc
}
