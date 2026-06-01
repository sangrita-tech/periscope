package output

type Writer interface {
	Write(content string) error
}
