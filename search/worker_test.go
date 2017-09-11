package search

import (
	"fmt"
	"testing"

	"os"

	"oschina.net/ContentSearch/filePool"
)

func Test(t *testing.T) {
	//wg := NewWorkerGroup(2)
	//wg.Debug()
	path, _ := os.Getwd()
	path = path + "/" + "a.doc"
	fs := filePool.FileSetNew()
	fs.Add(path)
	fs.Add(path)
	fs.Add(path)
	fs.Add(path)
	fs.Add(path)
	//fs.Close()
	//fs.Add("b.docx")
	//fs.Add("c.docx")
	//fs.Get()
	//fs.Get()
	//fmt.Printf("fs :%v\n", fs.IsEmpty())

	task := TaskInit("深度", 2)
	//task.Debug()
	task.Exec(fs)
	for _, s := range task.GetResult().v {
		fmt.Printf("result : %v\n", s)
	}

}
