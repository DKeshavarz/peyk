package codegen

type codeGenerator struct {
	groups    []int
	separator string
	charset   string
}

type Option func(*codeGenerator)

func WithGroups(groups []int) Option {
	return func(g *codeGenerator) {
		g.groups = groups
	}
}

func WithSeparator(separator string) Option {
	return func(g *codeGenerator) {
		g.separator = separator
	}
}

func WithCharset(charset string) Option {
	return func(g *codeGenerator) {
		g.charset = charset
	}
}

func WithLowercase() Option {
	return func(g *codeGenerator) {
		g.charset = "abcdefghijklmnopqrstuvwxyz"
	}
}

func WithUppercase() Option {
	return func(g *codeGenerator) {
		g.charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	}
}

func WithAlphanumeric() Option {
	return func(g *codeGenerator) {
		g.charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
}

func WithNumeric() Option {
	return func(g *codeGenerator) {
		g.charset = "0123456789"
	}
}

func NewCodeGenerator(options ...Option) CodeGenerator {
	defaultGen := &codeGenerator{
		groups:    []int{3, 4, 3},
		separator: "-",
		charset:   "abcdefghijklmnopqrstuvwxyz",
	}

	for _, option := range options {
		option(defaultGen)
	}

	return defaultGen
}
