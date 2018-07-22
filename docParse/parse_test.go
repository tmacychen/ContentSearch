package docParse

import (
	"fmt"
	"testing"
)

func TestParse(t *testing.T) {
	filename := "a.docx"
	parser := NewParser()
	parser.Init(filename)
	parser.Parse()
	fmt.Printf("parse.buf is %v\n", parser.buf)
}
