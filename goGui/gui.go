package goGui

import (
	"github.com/mattn/go-gtk/glib"
	"github.com/mattn/go-gtk/gtk"
)

// InitGui initializes the GUI
// TODO: Switch to a modern gui like fyne or walk
func InitGui() {
	gtk.Init(nil)
	window := gtk.NewWindow(gtk.WINDOW_TOPLEVEL)
	window.SetTitle("Job Aggregator")
	window.Connect("destroy", func(ctx *glib.CallbackContext) {
		gtk.MainQuit()
	}, "")
	window.SetPosition(gtk.WIN_POS_CENTER)
	window.SetDefaultSize(800, 600)

	vbox := gtk.NewVBox(false, 1)
	window.Add(vbox)

	// Create a new label
	label := gtk.NewLabel("Hello, Go-GTK!")
	vbox.Add(label)

	// Create a new button
	button := gtk.NewButtonWithLabel("Click me!")
	button.Connect("clicked", func() {
		label.SetText("Button was clicked!")
	})
	vbox.Add(button)

	window.ShowAll()
	gtk.Main()
}
