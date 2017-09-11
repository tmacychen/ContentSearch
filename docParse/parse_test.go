package docParse

import (
	"fmt"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	filename := "a.docx"
	var d DocType
	if strings.HasSuffix(filename, ".doc") {
		d = Word2003
	} else if strings.HasSuffix(filename, ".docx") {
		d = Word2007
	} else {
		d = UnknowType
	}
	parser := NewParser()
	parser.SetPath(d, filename)
	parser.Parse()
	fmt.Printf("parse.buf is %v\n", parser.buf)
}
