package worker

import (
	"math/rand"
	"time"

	"github.com/sirupsen/logrus"
)

func CreateWriteToChannelProcessorFunc(channel chan<- []byte) ProcessorFunc {
	return func(buffer []byte) {
		localBuf := make([]byte, len(buffer))
		copy(localBuf, buffer)
		channel <- localBuf
	}
}

func CreateWriteToConsoleProcessorFunc(direction string, logger *logrus.Logger) ProcessorFunc {
	return func(buffer []byte) {
		logger.WithFields(logrus.Fields{
			"direction":    direction,
			"content-size": len(buffer),
		}).Trace("DATA TRANSFER")
	}
}

func CreateMuddlingProcessorFunc() ProcessorFunc {
	return func(buffer []byte) {
		s1 := rand.NewSource(time.Now().UnixNano())
		r1 := rand.New(s1)

		muddledByteLocation := r1.Int31n((int32)(len(buffer)))

		buffer[muddledByteLocation] = byte(r1.Int31n(255))
	}
}
