package docParse

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	"oschina.net/ContentSearch/myerr"
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

func (p *Parser) SetPath(t DocType, path string) {
	switch t {
	case Word2003:
		p.parserFunc = docParse
	case Word2007:
		p.parserFunc = docxParse
	default:
		myerr.PrintErr("error in NewParser", errors.New("unkonw type of document"))
	}
	p.path = path
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

//测试用
func (p *Parser) ShowBuf() {
	fmt.Printf("buf :%v\n", p.buf)
}
