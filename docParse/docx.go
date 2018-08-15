package docParse

import (
	"archive/zip"
	"bufio"
	"bytes"
	"encoding/xml"
	"errors"
	"io"
	"os/exec"

	"github.com/apex/log"
)

func docxParse(filePath string, buf *bytes.Buffer) {
	var file *zip.File

	zipReader, e := zip.OpenReader(filePath)
	if e != nil {
		log.Fatalf("Open docx file [%v] error :%v\n", filePath, e)
	}
	defer zipReader.Close()

	for _, f := range zipReader.File {
		if f.Name == "word/document.xml" {
			file = f
		}
	}
	if file == nil {
		log.Fatalf("docx parser don't found document%v\n", errors.New("no found docx document"))
	}
	fileReader, e := file.Open()
	if e != nil {
		log.Fatalf("fileReader open %v\n", e)
	}
	textBuf := new(bytes.Buffer)

	//fInfo := file.FileInfo()
	//log.Infof("size (k):%v\n", fInfo.Size()/1024) // the unzip file size

	e = XMLToText(fileReader, textBuf, []string{"br", "p", "tab"}, []string{"instrText", "script"}, true)
	if e != nil {
		log.Fatalf("xml to text %v\n", e)

	}

	tr1 := exec.Command("tr", "-d", "\" \"")
	tr1.Stdin = bufio.NewReader(textBuf)
	tr1Reader, e := tr1.StdoutPipe()
	if e != nil {
		log.Fatalf("tr1: %v\n", e)
	}

	tr2 := exec.Command("tr", "-s", "\n")
	tr2.Stdin = tr1Reader
	tr2Reader, e := tr2.StdoutPipe()
	if e != nil {
		log.Fatalf("tr2: %v\n", e)
	}

	tr1.Start()
	tr2.Start()

	buf.ReadFrom(tr2Reader)
}

// XMLToText Convert XML to plain text given how to treat elements
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
			buf.Write([]byte(v))
		case xml.StartElement:
			for _, breakElement := range breaks {
				if v.Name.Local == breakElement {
					buf.Write([]byte("\n"))
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
