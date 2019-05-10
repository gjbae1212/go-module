package async_task

import "time"

type Option interface {
	apply(k *keeper)
}

type OptionFunc func(k *keeper)

func (f OptionFunc) apply(k *keeper) { f(k) }

func WithErrorHandlerOption(f ErrorHandler) OptionFunc {
	return func(k *keeper) {
		k.errorHandler = f
	}
}

func WithQueueSizeOption(size int) OptionFunc {
	return func(k *keeper) {
		k.queueSize = size
	}
}

func WithWorkerSizeOption(size int) OptionFunc {
	return func(k *keeper) {
		k.workerSize = size
	}
}

func WithTimeoutOption(timeout time.Duration) OptionFunc {
	return func(k *keeper) {
		k.timeout = timeout
	}
}
