package mainWidget

import (
	"io/ioutil"
	"log"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"oschina.net/ContentSearch/filePool"
	"oschina.net/ContentSearch/myerr"
	sech "oschina.net/ContentSearch/search"

	"github.com/gotk3/gotk3/gtk"
)

var numCPU = runtime.NumCPU()

func MainWidget(win *gtk.Window) *gtk.Widget {
	runtime.GOMAXPROCS(numCPU)
	grid, err := gtk.GridNew()
	myerr.PrintErr("Unable to create grid:", err)

	grid.SetOrientation(gtk.ORIENTATION_VERTICAL)
	grid.SetBorderWidth(5)

	grid1, err := gtk.GridNew()
	myerr.PrintErr("Unable to create grid:", err)

	grid1.SetOrientation(gtk.ORIENTATION_HORIZONTAL)

	grid2, err := gtk.GridNew()
	myerr.PrintErr("Unable to create grid:", err)

	grid2.SetOrientation(gtk.ORIENTATION_HORIZONTAL)

	grid.Add(grid1)
	grid.AttachNextTo(grid2, grid1, gtk.POS_BOTTOM, 1, 1)

	treeView := TreeViewInit()
	InitColumn(treeView)

	label, err := gtk.LabelNew("")
	myerr.PrintErr("Unable to create label:", err)

	entry, err := gtk.EntryNew()
	myerr.PrintErr("Unable to create entry:", err)

	entry.SetHExpand(true)

	entry.Connect("activate", func() {
		s, _ := entry.GetText()
		label.SetText(s)
	})

	btnfile, err := gtk.FileChooserButtonNew("选择一个文件夹", gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)
	myerr.PrintErr("Open a directory error :%v\n", err)

	u, err := user.Current()
	myerr.PrintErr("get current user err:%v\n", err)

	var fs *filePool.FileSet
	var task *sech.Task

	btnfile.SetCurrentFolder(u.HomeDir)
	btnfile.Connect("file-set", func() {
		dir := btnfile.GetFilename()
		label.SetText(dir)
		listStore := ListStoreInit(treeView)
		fs = filePool.FileSetNew()
		if treeView.GetNColumns() < 3 {
			AddTimeCol(treeView)
		}
		go func() {
			openDir(dir, listStore, fs)
			fs.Close()
		}()
		//第一次搜索后，再次选择文件夹，列出所有文件后，
		//需要清空上次结果，否则打开按钮会打开上次结果的文档
		if task != nil {
			task.ClearResult()
		}
	})

	btnSearch, err := gtk.ButtonNewWithLabel("搜索")
	myerr.PrintErr("Unable to create button:", err)

	btnSearch.Connect("clicked", func() {
		text, _ := entry.GetText()
		if text == "" {
			ShowMessage(win, "输入的文本为空")
			return
		}
		listStore := ListStoreInit(treeView)
		RemoveTimeCol(treeView)
		task = search(win, text, listStore, fs)
	})

	btnOpen, err := gtk.ButtonNewWithLabel("打开")
	myerr.PrintErr("Unable to create button:", err)

	btnOpen.Connect("clicked", func() {
		if task == nil {
			ShowMessage(win, "你需要输入搜索的内容，然后点击搜索按钮")
			log.Println("the task is not exist")
			return
		}
		if task.GetResult().ItemLen() == 0 {
			log.Println("the item is unavailable")
			return
		}
		err := exec.Command("xdg-open", task.GetResult().GetPath(selectItem)).Run()
		myerr.PrintErr("xdg-open err:%s", err)
	})

	grid1.Add(btnfile)
	grid1.SetColumnSpacing(2)
	grid1.AttachNextTo(entry, btnfile, gtk.POS_RIGHT, 5, 5)
	grid1.AttachNextTo(btnSearch, entry, gtk.POS_RIGHT, 5, 5)
	grid1.AttachNextTo(btnOpen, btnSearch, gtk.POS_RIGHT, 5, 5)
	grid2.Add(label)

	scrolledWindow, err := gtk.ScrolledWindowNew(nil, nil)
	myerr.PrintErr("Unable create scrolledWindow :", err)

	scrolledWindow.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)
	scrolledWindow.SetBorderWidth(5)
	scrolledWindow.SetHExpand(true)
	scrolledWindow.SetVExpand(true)
	scrolledWindow.Add(treeView)
	grid.AttachNextTo(scrolledWindow, grid2, gtk.POS_BOTTOM, 1, 1)

	return &grid.Container.Widget
}

//递归调用搜索所有的目录
func openDir(dirPath string, listStore *gtk.ListStore, fs *filePool.FileSet) {
	//println("directory name :", dirPath)

	dir, err := ioutil.ReadDir(dirPath)
	myerr.PrintErr("open directory:"+dirPath+"error:", err)

	//遍历目录下的内容，获取文件详情，同os.Stat(filename)获取的信息
	for _, info := range dir {
		n := info.Name()
		p := dirPath + "/" + n
		if info.IsDir() {
			openDir(p, listStore, fs)
			continue
		}
		if !isReadableFile(n) {
			continue
		} //文件名
		//info.Mode()    //文件权限
		//info.Size()    //文件大小
		//info.Sys()     //系统信息
		fs.Add(p)
		log.Println("fs add file " + p)
		AddRow(listStore, n, info.ModTime().Format("2006-01-02 15:04:05"), info.Size()/1000)
	}

}

func isReadableFile(name string) bool {
	if strings.HasSuffix(name, ".doc") {
		return true
	}
	if strings.HasSuffix(name, ".docx") {
		return true
	}
	return false

}

func ShowMessage(win *gtk.Window, text string) {
	dialog := gtk.MessageDialogNew(win, gtk.DIALOG_MODAL,
		gtk.MESSAGE_INFO, gtk.BUTTONS_CLOSE, "%s", text)
	dialog.SetTitle("注意了！！")
	dialog.SetHExpand(true)
	win.Add(dialog)
	dialog.SetDefaultSize(100, 50)
	b, _ := dialog.GetWidgetForResponse(gtk.RESPONSE_CLOSE)
	b.Connect("clicked", func() {
		dialog.Destroy()
	})
	dialog.Show()
}
func search(win *gtk.Window, text string, ls *gtk.ListStore, fs *filePool.FileSet) *sech.Task {

	if fs == nil {
		log.Println("fs is nil")
		return nil
	}

	//暂时支持单个词搜索
	//only search one word at present
	t := sech.TaskInit(sech.Key(text), numCPU-1)
	t.Exec(fs)
	r := t.GetResult()
	for i := 0; i <= r.ItemLen(); i++ {
		AddResultRow(ls, r.GetName(i), r.GetContent(i))
	}
	return t
}
