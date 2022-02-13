package config

import (
	"time"

	"go.mongodb.org/mongo-driver/mongo/options"
)

type Options struct {
	MaxLogs         int64
	MaxSize         int64
	MaxBufferSize   int64
	FlushEvery      time.Duration
	AppName         string
	Tags            []string
	Database        string
	DatabaseOptions *options.ClientOptions
	Verbose         bool
}

type Option func(opts *Options)

func WithVerbose(enabled bool) Option {
	return func(opts *Options) {
		opts.Verbose = enabled
	}
}

func WithDatabaseOptions(opts *options.ClientOptions) Option {
	return func(o *Options) {
		o.DatabaseOptions = opts
	}
}

func WithMaxBufferSize(size int64) Option {
	return func(opts *Options) {
		opts.MaxBufferSize = size
	}
}

func WithTags(tags ...string) Option {
	return func(opts *Options) {
		opts.Tags = tags
	}
}

func WithName(name string) Option {
	return func(opts *Options) {
		opts.AppName = name
	}
}

func WithMaxLogs(max int64) Option {
	return func(opts *Options) {
		opts.MaxSize = max
	}
}

func WithMaxSize(bytes int64) Option {
	return func(opts *Options) {
		opts.MaxSize = bytes
	}
}

func WithFlushEvery(dur time.Duration) Option {
	return func(opts *Options) {
		opts.FlushEvery = dur
	}
}
