package dnet

import "time"

type Option func(opt *Options)

// loadOptions returns an initialized *Options with options
func loadOptions(options ...Option) *Options {
	opts := new(Options)
	for _, option := range options {
		option(opts)
	}
	return opts
}

// copyOption returns an *Options with options
func copyOption(src *Options, opts ...Option) *Options {
	for _, option := range opts {
		option(src)
	}
	return src
}

// Options contains all options which will be applied when instantiating a session.
type Options struct {
	// when the send channel is full, BlockSend will block if it is true.
	// or it return queue full error code, if BlockSend is false.
	// it default is false
	BlockSend bool

	// capacity of the send channel. default net.defSendChannelSize
	SendChannelSize int

	// the deadline for read
	ReadTimeout time.Duration

	// the deadline for write
	WriteTimeout time.Duration

	// session will call the MsgCallback,if it has a message
	MsgCallback func(session Session, message interface{})

	// session will call the ErrorCallback,if it has a error
	ErrorCallback func(session Session, err error)

	// session will call the CloseCallback,if it is closed
	CloseCallback func(session Session, reason error)

	// encoder and decoder
	Codec Codec
}

// WithOptions accepts the whole options config.
func WithOptions(option *Options) Option {
	return func(opt *Options) {
		opt = option
	}
}

// WithBlockSend indicates whether it should block when the send channel full.
func WithBlockSend(bs bool) Option {
	return func(opt *Options) {
		opt.BlockSend = bs
	}
}

// WithSendChannelSize sets capacity of the send channel .
func WithSendChannelSize(size int) Option {
	return func(opt *Options) {
		opt.SendChannelSize = size
	}
}

// WithMessageCallback sets message callback.
func WithMessageCallback(msgCb func(session Session, message interface{})) Option {
	return func(opt *Options) {
		opt.MsgCallback = msgCb
	}
}

// WithErrorCallback sets error callback.
func WithErrorCallback(errCb func(session Session, err error)) Option {
	return func(opt *Options) {
		opt.ErrorCallback = errCb
	}
}

// WithTimeout sets the deadline of read/write.
func WithTimeout(readTimeout, writeTimeout time.Duration) Option {
	return func(opt *Options) {
		opt.ReadTimeout = readTimeout
		opt.WriteTimeout = writeTimeout
	}
}

// WithCodec sets codec.
func WithCodec(codec Codec) Option {
	return func(opt *Options) {
		opt.Codec = codec
	}
}

// WithCloseCallback sets close callback.
func WithCloseCallback(closeCallback func(session Session, reason error)) Option {
	return func(opt *Options) {
		opt.CloseCallback = closeCallback
	}
}
