package docParse

import (
	"bytes"
	"errors"
	"strings"
)

// DocType is an Int type
type DocType int

const (
	UnknowType DocType = iota //int type 0
	Word2003                  // doc type is 1
	Word2007                  // docx type is 2
)

type Parser struct {
	path       string
	buf        *bytes.Buffer
	parserFunc func(string, *bytes.Buffer)
}

func NewParser() *Parser {
	p := new(Parser)
	p.buf = new(bytes.Buffer)
	return p
}

func (p *Parser) Init(path string) error {
	if strings.HasSuffix(path, ".doc") {
		p.parserFunc = docParse
	} else if strings.HasSuffix(path, ".docx") {
		p.parserFunc = docxParse
	} else {
		return errors.New("unkonw type of document")
	}
	p.path = path
	return nil
}
func (p *Parser) Path() string {
	return p.path
}
func (p *Parser) FileName() string {
	s := strings.Split(p.path, "/")
	return s[len(s)-1]
}
func (p *Parser) Parse() {
	p.parserFunc(p.path, p.buf)
}
func (p *Parser) ClearBuf() {
	p.buf.Reset()
}
func (p *Parser) GetBuf() *bytes.Buffer {
	return p.buf
}

//ShowBuf
//func (p *Parser) ShowBuf() {
//	fmt.Printf("buf :%v\n", p.buf)
//}
//
