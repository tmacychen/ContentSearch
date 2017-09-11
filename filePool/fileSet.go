package filePool

type FileSet struct {
	filePathChan chan string
}

func FileSetNew() *FileSet {
	fs := new(FileSet)
	fs.filePathChan = make(chan string, 100)
	return fs
}

//add 操作可能不需要加锁
func (fs *FileSet) Add(s string) {
	fs.filePathChan <- s
}

//Get是多个worker同时获取任务，获取文件路径为互斥的内容
func (fs *FileSet) Get() string {
	return <-fs.filePathChan
}

func (fs *FileSet) Length() int {
	return len(fs.filePathChan)
}
func (fs *FileSet) Close() {
	close(fs.filePathChan)
}
