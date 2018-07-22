package search

import (
	"fmt"
	"strings"
	"unicode"

	"github.com/apex/log"
	"github.com/tmacychen/ContentSearch/docParse"
)

type Worker struct {
	parser    *docParse.Parser
	isWorking bool
}

func NewWorker() *Worker {
	w := &Worker{}
	w.parser = docParse.NewParser()
	w.isWorking = false
	return w
}

const BOUND = 60

//Do 工作者会处理文档类型，选择对应的解析器将其转化为文本文档
//
func (w *Worker) Do(filePath string, key Key, res *Result) {

	if err := w.parser.Init(filePath); err != nil {
		log.Errorf("work parse file errr :%v\n", err)
		return
	}
	w.parser.Parse()
	//	fmt.Printf("search :%v\n", filePath)
	search(w.parser, key, res)
	w.parser.ClearBuf()
}

func (w *Worker) iamBusy(s bool) {
	w.isWorking = s
}

//IsBusy shows that whether the woker is busy
func (w *Worker) IsBusy() bool {
	return w.isWorking
}

//搜索给定关键字，词在文中的位置，并将上下文作为返回值加入到res中
func search(p *docParse.Parser, k Key, res *Result) {

	l1 := []byte{13, 10} // 在wps中，回车与换行为CR和LF符号，不是"\n"
	l2 := []byte{10}

	space := []byte{32}

	b := p.GetBuf().String()
	lenOfKey := len(k)

	for {
		n := strings.Index(b, string(k))
		var r string
		//if find the key word, put it to result slice
		//if not ,go return
		if n != -1 {
			a := n + lenOfKey

			if a+BOUND < len(b) {
				if n-BOUND < 0 { // 前边界
					//为了防止乱码，处理边界，如果不是整除3，删除前1或2个字节。
					c := n % 3
					r = b[c : a+BOUND]
				} else {
					r = b[n-BOUND : a+BOUND]
				}
				r = cleanUnreadableSymbol(r)
				b = b[a+BOUND:]
			} else {
				if n-BOUND > 0 {
					b = b[n-BOUND:]
				}
				r = cleanUnreadableSymbol(b)
				r = strings.Replace(r, string(l1), string(space), -1) // 清除要显示内容中的回车与换行,换成空格
				r = strings.Replace(r, string(l2), string(space), -1) // 清除要显示内容中的换行,换成空格
				//log.Printf("the Conntent is %v\n", string(r))
				res.AddOneItem(p.Path(), p.FileName(), r)
				return
			}
			r = strings.Replace(r, string(l1), string(space), -1) // 清除要显示内容中的回车与换行,换成空格
			r = strings.Replace(r, string(l2), string(space), -1) // 清除要显示内容中的换行,换成空格
			//log.Printf("the Conntent is %v\n", string(r))
			res.AddOneItem(p.Path(), p.FileName(), r)

		} else {
			return
		}
	}
}

func cleanHead(r []rune) string {
	for i := 0; i < 5; i++ {
		if unicode.Is(unicode.Scripts["Han"], r[i]) {
			r = r[i:]
			break
		}
	}
	fmt.Println("clean Head")
	return string(r)
}

func cleanUnreadableSymbol(s string) string {
	r := []rune(s)
	for i := 0; i < 5; i++ {
		if unicode.IsLetter(r[i]) || unicode.Is(unicode.Scripts["Han"], r[i]) {
			r = r[i:]
			break
		}
	}
	l := len(r) - 1
	for i := l; i > l-3; i-- {
		if unicode.IsLetter(r[i]) || unicode.Is(unicode.Scripts["Han"], r[i]) {
			return string(r[:i+1])
		}
	}
	return string(r[:l-3])
}
