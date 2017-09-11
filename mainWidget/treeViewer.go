package mainWidget

import (
	"log"
	"strconv"

	"oschina.net/ContentSearch/myerr"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

// IDs to access the tree view columns by
const (
	COLUMN_FILENAME = iota
	COLUMN_CONTENT
	COLUMN_TIME
)

var selectItem int
var nameCol *gtk.TreeViewColumn
var timeCol *gtk.TreeViewColumn
var contentCol *gtk.TreeViewColumn

// Creating a list store. This is what holds the data that will be shown on our tree view.
func ListStoreInit(treeView *gtk.TreeView) *gtk.ListStore {
	listStore, err := gtk.ListStoreNew(glib.TYPE_STRING, glib.TYPE_STRING, glib.TYPE_STRING)
	if err != nil {
		log.Fatal("Unable to create list store:", err)
	}
	treeView.SetModel(listStore)
	return listStore
}

// Creates a tree view and the list store that holds its data
func TreeViewInit() *gtk.TreeView {
	treeView, err := gtk.TreeViewNew()
	if err != nil {
		log.Fatal("Unable to create tree view :", err)
	}

	treeView.SetActivateOnSingleClick(true)
	selection, err := treeView.GetSelection()
	if err != nil {
		log.Fatal("Could not get tree selection object.")
	}
	selection.SetMode(gtk.SELECTION_SINGLE)
	selection.Connect("changed", oneRowSelected)
	return treeView
}

func InitColumn(tv *gtk.TreeView) {

	nameCol = createTextColumn("文件名称", COLUMN_FILENAME, 1000)
	timeCol = createTextColumn("日期", COLUMN_TIME, 300)
	contentCol = createTextColumn("大小(KB)", COLUMN_CONTENT, 300)

	tv.AppendColumn(nameCol)
	tv.AppendColumn(timeCol)
	tv.AppendColumn(contentCol)
}

func RemoveTimeCol(tv *gtk.TreeView) {
	tv.RemoveColumn(timeCol)
	tv.RemoveColumn(contentCol)
	contentCol = createTextColumn("内容", COLUMN_CONTENT, 1000)
	tv.AppendColumn(contentCol)
}
func AddTimeCol(tv *gtk.TreeView) {
	tv.RemoveColumn(contentCol)
	timeCol = createTextColumn("日期", COLUMN_TIME, 150)
	contentCol = createTextColumn("大小(KB)", COLUMN_CONTENT, 1000)
	tv.AppendColumn(timeCol)
	tv.AppendColumn(contentCol)
}

func oneRowSelected(sel *gtk.TreeSelection) {
	var iter *gtk.TreeIter
	var model gtk.ITreeModel
	var ok bool

	model, iter, ok = sel.GetSelected()
	if ok {
		tpath, err := model.(*gtk.TreeModel).GetPath(iter)
		myerr.PrintErr("treeSelectionChangedCB: Could not get path from model: %s\n", err)

		selectItem, err = strconv.Atoi(tpath.String())
		myerr.PrintErr("strconv from string to int err :%s", err)

	}
}

// Add a column to the tree view (during the initialization of the tree view)
func createTextColumn(title string, id, width int) *gtk.TreeViewColumn {
	cellRenderer, err := gtk.CellRendererTextNew()
	if err != nil {
		log.Fatal("Unable to create text cell renderer:", err)
	}
	column, err := gtk.TreeViewColumnNewWithAttribute(title, cellRenderer, "text", id)
	if err != nil {
		log.Fatal("Unable to create cell column:", err)
	}
	column.SetClickable(true)
	column.SetSpacing(5)
	column.SetMaxWidth(width)

	column.SetResizable(true) // 栏可以拖拽改变宽度
	//显示每一栏上的排序箭头标志
	column.SetSortIndicator(true)
	column.SetSortColumnID(1)

	return column
}

func AddRow(listStore *gtk.ListStore, fileName, time string, size int64) {
	// Get an iterator for a new row at the end of the list store
	iter := listStore.Append()

	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter,
		[]int{COLUMN_FILENAME, COLUMN_TIME, COLUMN_CONTENT},
		[]interface{}{fileName, time, size})
	if err != nil {
		log.Fatal("Unable to add row:", err)
	}
}
func AddResultRow(listStore *gtk.ListStore, fileName, content string) {
	// Get an iterator for a new row at the end of the list store
	iter := listStore.Append()

	// Set the contents of the list store row that the iterator represents
	err := listStore.Set(iter,
		[]int{COLUMN_FILENAME, COLUMN_CONTENT},
		[]interface{}{fileName, content})
	myerr.PrintErr("Unable to add row:", err)
}
