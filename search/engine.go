package search

import (
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/tmacychen/ContentSearch/filePool"
)

//Key 要检索的关键字符串
type Key string

type item struct {
	absPath  string // 绝对路径
	fileName string //文件名称
	content  string //内容
}

//Result 搜索结果
type Result struct {
	v      []item
	locker *sync.Mutex
}

//AddOneItem 增加一个结果，并发安全！
func (r *Result) AddOneItem(path, name, content string) {
	i := item{path, name, content}
	r.locker.Lock()
	r.v = append(r.v, i)
	r.locker.Unlock()
}

//GetPath get the item's file path
func (r *Result) GetPath(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].absPath
	} else {
		return ""
	}
}

//GetName get the item's file's name
func (r *Result) GetName(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].fileName
	} else {
		return ""
	}
}

//GetContent get the item's content
func (r *Result) GetContent(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].content
	} else {
		return ""
	}
}

//ItemLen return the length of items
func (r *Result) ItemLen() int {
	return len(r.v)
}

//Task 一次任务，描述了一次搜索的执行.It's what the workers do.
type Task struct {
	key       Key
	res       *Result
	workerNum int
	end       bool
	workers   []*Worker //worker set
}

//TaskInit 初始化任务
//参数 key：需要搜索的关键字 n：并发数量
//返回值：初始化完成的Task指针
//通常n与处理器个数相同
func TaskInit(k Key, n int) *Task {
	if n <= 0 {
		log.Fatalf("worker group's number < 0")
		return nil
	}
	t := &Task{
		key:       k,
		workerNum: n,
		end:       false,
		res: &Result{
			locker: new(sync.Mutex),
		},
	}

	for n > 0 {
		w := NewWorker()
		t.workers = append(t.workers, w)
		n--
	}
	return t
}

//GetKey 获取当前任务的关键字
func (t *Task) GetKey() Key {
	return t.key
}

//SetEnd 设置当前任务准备结束。当文件集合读取所有文件后，会设置
//task的状态为end,task会执行所有任务后结束
func (t *Task) SetEnd(s bool) {
	t.end = s
}

//GetResult 获取此任务的全部结果
func (t *Task) GetResult() *Result {
	return t.res
}

//ClearResult 清空此任务(用于保留任务信息，清楚历史信息)
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
			t.workers[i].SetBusy(true)
			return t.workers[i]
		}
	}
	//log.Debug("get worker nil")
	return nil
}

//Exec 执行任务
//参数：文件集合*FileSet
//无返回值
//并发执行对文件的解析，并对关键字查找.等待所有子线程结束后，退出
func (t *Task) Exec(fs *filePool.FileSet) {
	var wait sync.WaitGroup
	for !t.end || fs.Length() > 0 {
		w := t.getWorker()
		if w != nil {
			log.Debugf("get worker fs.len:%v\n", fs.Length())
			wait.Add(1)
			go func() {
				s := fs.Get()
				log.Debugf("end :%v \t len :%d \nExec fs.Get:%v\n", t.end, fs.Length(), s)
				if s != "" {
					w.Do(s, t.key, t.res)
				}
				w.SetBusy(false)
				wait.Done()
			}()
		} else {
			time.Sleep(time.Millisecond * 1)
		}

	}
	wait.Wait()
}

//Debug 调试用
func (t *Task) Debug() {
	log.Debugf("t.key:%v\n", t.GetKey())
}
