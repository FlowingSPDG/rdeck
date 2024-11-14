package output

type VMixOutput interface {
	SendFunction(name, query string) error
	// TODO: more!
}
