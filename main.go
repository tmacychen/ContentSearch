package main

import (
	"log"

	"github.com/gotk3/gotk3/gtk"
	"oschina.net/ContentSearch/mainWidget"
)

var (
	Version   string
	BuildTime string
)

func main() {
	gtk.Init(nil)

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	win.SetTitle("Content Search")
	win.Connect("destroy", func() {
		gtk.MainQuit()
	})
	// Set the default window size.
	win.SetDefaultSize(1000, 500)
	win.Add(mainWidget.MainWidget(win))

	win.ShowAll()
	gtk.Main()

}
