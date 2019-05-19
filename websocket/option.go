package websocket

type Option interface {
	apply(bk *breaker)
}

type OptionFunc func(bk *breaker)

func (f OptionFunc) apply(bk *breaker) { f(bk) }

func WithErrorHandlerOption(f ErrorHandler) OptionFunc {
	return func(bk *breaker) {
		bk.errorHandler = f
	}
}

func WithMaxReadLimit(length int64) OptionFunc {
	return func(bk *breaker) {
		bk.maxReadLimit = length
	}
}

func WithMaxMessagePoolLength(length int64) OptionFunc {
	return func(bk *breaker) {
		bk.broadcast = make(chan Message, length)
	}
}
