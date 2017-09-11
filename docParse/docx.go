package docParse

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"os/exec"

	"oschina.net/ContentSearch/myerr"
)

func docxParse(filePath string, buf *bytes.Buffer) {
	var file *zip.File
	f, e := os.OpenFile(filePath, os.O_RDONLY, 0440)
	myerr.PrintErr("OpenFile error :", e)
	defer f.Close()

	zipReader, e := zip.OpenReader(filePath)
	myerr.PrintErr("open docx file error", e)
	defer zipReader.Close()
	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			file = f
		}
	}
	if file == nil {
		myerr.PrintErr("docx parser don't found document", errors.New("no found docx document"))
	}
	fileReader, e := file.Open()
	myerr.PrintErr("fileReader open error ", e)
	textBuf := new(bytes.Buffer)
	e = XMLToText(fileReader, textBuf, []string{"br", "p", "tab"}, []string{"instrText", "script"}, true)
	myerr.PrintErr("xml to text", e)

	tr1 := exec.Command("tr", "-d", "\" \"")
	tr1.Stdin = bufio.NewReader(textBuf)
	tr1Reader, e := tr1.StdoutPipe()
	myerr.PrintErr("tr1:", e)

	tr2 := exec.Command("tr", "-s", "\n")
	tr2.Stdin = tr1Reader
	tr2Reader, e := tr2.StdoutPipe()
	myerr.PrintErr("tr2:", e)

	tr1.Start()
	tr2.Start()

	buf.ReadFrom(tr2Reader)
}

// Convert XML to plain text given how to treat elements
func XMLToText(r io.Reader, buf *bytes.Buffer, breaks []string, skip []string, strict bool) error {

	dec := xml.NewDecoder(r)
	dec.Strict = strict
	for {
		t, err := dec.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		switch v := t.(type) {
		case xml.CharData:
			buf.WriteString(string(v))
		case xml.StartElement:
			for _, breakElement := range breaks {
				if v.Name.Local == breakElement {
					buf.WriteString("\n")
				}
			}
			for _, skipElement := range skip {
				if v.Name.Local == skipElement {
					depth := 1
					for {
						t, err := dec.Token()
						if err != nil {
							// An io.EOF here is actually an error.
							return err
						}
						switch t.(type) {
						case xml.StartElement:
							depth++
						case xml.EndElement:
							depth--
						}

						if depth == 0 {
							break
						}
					}
				}
			}
		}
	}
	return nil
}
