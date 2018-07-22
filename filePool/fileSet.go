package filePool

//FileSet save the file path in channels
type FileSet struct {
	fileChannels []chan string
	top          int
}

const bufferSize = 10

//FileSetNew create a new file set
func FileSetNew() *FileSet {
	fs := &FileSet{
		top: 0,
	}
	c := make(chan string, bufferSize)
	fs.fileChannels = append(fs.fileChannels, c)
	return fs
}

//Add add a file path to fileSet
func (fs *FileSet) Add(s string) {
	//找一个空channel，传入string
	if len(fs.fileChannels[fs.top]) == bufferSize {
		fs.fileChannels = append(fs.fileChannels, make(chan string, bufferSize))
		fs.top++
	}
	fs.fileChannels[fs.top] <- s
}

//bug:当work少于4个时，会堵塞在get中没有退出

//Get the works get task from fileSet
//是多个worker同时获取任务，获取文件路径为互斥的内容
func (fs *FileSet) Get() string {
	for _, c := range fs.fileChannels {
		select {
		case s := <-c:
			return s
		default:
		}
	}
	return ""
}

//Length return the number of file in the current fileSet
func (fs *FileSet) Length() int {
	sum := 0
	for _, c := range fs.fileChannels {
		sum += len(c)
	}
	return sum
}
