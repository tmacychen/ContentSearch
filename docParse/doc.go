package docParse

import (
	"bytes"
	//"fmt"
	"github.com/apex/log"
	"io/ioutil"
	"os"
	"os/exec"
)

func docParse(filePath string, buf *bytes.Buffer) {

	//Save output to a file
	outputFile, e := ioutil.TempFile("/tmp", "doc-")
	if e != nil {
		log.Fatalf("TempFile Out:%v\n", e)
	}

	//	fmt.Println("file name ",f.Name()," output File name ",outputFile.Name())
	e = exec.Command("wvText", filePath, outputFile.Name()).Run()
	if e != nil {
		log.Infof("wvText parser [%v] error :%v\n", filePath, e) //文档解析失败返回，不退出本次程序
		os.Remove(outputFile.Name())
		return
	}

	cat := exec.Command("cat", outputFile.Name())
	catReader, e := cat.StdoutPipe()
	if e != nil {
		log.Fatalf("cat %v\n", e)
	}
	//采用tr命令，清楚多余空格
	tr1 := exec.Command("tr", "-d", "\" \"")
	tr1.Stdin = catReader
	tr1Reader, e := tr1.StdoutPipe()
	if e != nil {
		log.Fatalf("tar -d %v\n", e)
	}

	tr2 := exec.Command("tr", "-s", "\n")
	tr2.Stdin = tr1Reader
	tr2Reader, e := tr2.StdoutPipe()
	if e != nil {
		log.Fatalf("tar -s %v\n", e)
	}

	cat.Start()
	tr1.Start()
	tr2.Start()

	_, e = buf.ReadFrom(tr2Reader)
	if e != nil {
		log.Fatalf("wvText: %v\n", e)
	}
	os.Remove(outputFile.Name())
}
