package search

import (
	"fmt"
	"log"
	"sync"
	"time"

	"oschina.net/ContentSearch/filePool"
)

type Key string

type item struct {
	absPath  string
	fileName string
	content  string
}
type Result struct {
	v      []item
	locker *sync.Mutex
}

func (r *Result) AddOneItem(path, name, content string) {
	i := item{path, name, content}
	r.locker.Lock()
	r.v = append(r.v, i)
	r.locker.Unlock()
}
func (r *Result) GetPath(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].absPath
	} else {
		return ""
	}
}
func (r *Result) GetName(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].fileName
	} else {
		return ""
	}
}
func (r *Result) GetContent(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].content
	} else {
		return ""
	}
}
func (r *Result) ItemLen() int {
	return len(r.v)
}

type Task struct {
	key       Key
	res       *Result
	workerNum int
	workers   []*Worker //worker set
}

//初始化任务
//参数 key：需要搜索的关键字 n：并发数量
//返回值：初始化完成的Task指针
//通常n与处理器个数相同
func TaskInit(key Key, n int) *Task {
	if n <= 0 {
		log.Fatalf("worker group's number < 0")
		return nil
	}
	t := new(Task)
	t.key = key
	t.workerNum = n
	t.res = new(Result)
	t.res.locker = new(sync.Mutex)

	for n > 0 {
		w := NewWorker()
		t.workers = append(t.workers, w)
		n--
	}
	return t
}

//获取当前任务的关键字
func (t *Task) GetKey() Key {
	return t.key
}

//获取此任务的全部结果
func (t *Task) GetResult() *Result {
	return t.res
}

func (t *Task) ClearResult() {
	t.res = new(Result)
	t.res.locker = new(sync.Mutex)
}

//获取此任务的可用worker的数量
//
//func (t *Task) GetWorkNum() int {
//	return len(t.workers)
//}

//获取一个设置为busy的worker
func (t *Task) getWorker() *Worker {
	for i := 0; i < len(t.workers); i++ {
		if !t.workers[i].IsBusy() {
			t.workers[i].iamBusy(true)
			return t.workers[i]
		}
	}
	return nil
}

//执行任务
//参数：文件集合*FileSet
//无返回值
//并发执行对文件的解析，并对关键字查找.等待所有子线程结束后，退出

func (t *Task) Exec(fs *filePool.FileSet) {
	var wait sync.WaitGroup

	for fs.Length() > 0 {
		w := t.getWorker()
		if w != nil {
			wait.Add(1)
			go func() {
				s := fs.Get()
				if s != "" {
					log.Printf("Exec fs.Get:%v\n", s)
					w.Do(s, t.key, t.res)
				}
				w.iamBusy(false)
				wait.Done()
			}()
		} else {
			time.Sleep(time.Millisecond * 1)
		}

	}
	wait.Wait()
}
func (t *Task) Debug() {
	fmt.Println("t.key:", t.GetKey())
}
