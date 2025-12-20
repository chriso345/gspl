package lang

// WithSourceURI sets the source URI on ParseOptions.
func WithSourceURI(uri string) ParseOption {
	return func(o *ParseOptions) { o.SourceURI = uri }
}

// WithMaxBytes sets a maximum number of bytes the parser will read.
func WithMaxBytes(n int64) ParseOption {
	return func(o *ParseOptions) { o.MaxBytes = n }
}
