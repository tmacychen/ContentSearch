package search

import (
	"fmt"
	//	"path"
	"strings"
	"sync"
	"time"

	"github.com/apex/log"
	"github.com/apex/log/handlers/text"
	"github.com/tmacychen/ContentSearch/filePool"
)

//Key 要检索的关键词
type Key string

type item struct {
	path    string   // 绝对路径
	content []string //内容
}

//Result 搜索结果
type Result struct {
	v      []item
	locker *sync.Mutex
}

//AddOneItem 增加一个结果，并发安全！
func (r *Result) AddOneItem(path, content string) {
	r.locker.Lock()
	defer r.locker.Unlock()

	for i := 0; i < len(r.v); i++ {
		if r.v[i].path == path {
			r.v[i].content = append(r.v[i].content, content)
			return
		}
	}
	i := item{}
	i.path = path
	i.content = append(i.content, content)
	r.v = append(r.v, i)
}

//GetPath get the item's file path
func (r *Result) GetPath(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].path
	} else {
		return ""
	}
}

//GetName get the item's file's name
func (r *Result) GetName(i int) string {
	if i >= 0 && i < len(r.v) {
		return r.v[i].path
	} else {
		return ""
	}
}

//GetContent get the item's content
// func (r *Result) GetContent(i int) string {
// 	if i >= 0 && i < len(r.v) {
// 		return r.v[i].content
// 	} else {
// 		return ""
// 	}
// }

//ItemLen return the length of items
func (r *Result) ItemLen() int {
	return len(r.v)
}

//Task 一次任务，描述了一次搜索的执行.It's what the workers do.
//key 表示搜索的关键词
//res 表示搜索结果
//workerNum 表示需要工人数量
//end 表示此次任务是否结束
//workers 工人集合

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
				fmt.Printf(">")
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

// ShowResult 显示结果
func (t *Task) ShowResult() {
	for i := 0; i < t.res.ItemLen(); i++ {
		fmt.Printf("\033[%dm%6s\n", text.Colors[log.InfoLevel], t.res.GetName(i))
		for j := 0; j < len(t.res.v[i].content); j++ {
			s := strings.Split(t.res.v[i].content[j], string(t.key))
			fmt.Printf("\t\033[0m%s\033[%dm%s\033[0m%s\n", s[0], text.Colors[log.ErrorLevel], t.key, s[1])
		}
	}
}
