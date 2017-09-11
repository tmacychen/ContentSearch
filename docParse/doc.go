package docParse

import (
	"bytes"

	"oschina.net/ContentSearch/myerr"
	//"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

func docParse(filePath string, buf *bytes.Buffer) {

	//Save output to a file
	outputFile, e := ioutil.TempFile("/tmp", "doc-")
	myerr.PrintErr("TempFile Out:", e)

	//	fmt.Println("file name ",f.Name()," output File name ",outputFile.Name())
	e = exec.Command("wvText", filePath, outputFile.Name()).Run()
	if e != nil {
		log.Printf("wvText error output :%v\n", e) //文档解析失败返回，不退出本次程序
		os.Remove(outputFile.Name())
		return
	}

	cat := exec.Command("cat", outputFile.Name())
	catReader, e := cat.StdoutPipe()
	myerr.PrintErr("cat :", e)
	//采用tr命令，清楚多余空格
	tr1 := exec.Command("tr", "-d", "\" \"")
	tr1.Stdin = catReader
	tr1Reader, e := tr1.StdoutPipe()
	myerr.PrintErr("cat :", e)

	tr2 := exec.Command("tr", "-s", "\n")
	tr2.Stdin = tr1Reader
	tr2Reader, e := tr2.StdoutPipe()
	myerr.PrintErr("cat :", e)

	cat.Start()
	tr1.Start()
	tr2.Start()

	_, e = buf.ReadFrom(tr2Reader)
	myerr.PrintErr("wvText:", e)
	os.Remove(outputFile.Name())
}
