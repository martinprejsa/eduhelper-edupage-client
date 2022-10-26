package gtk

import (
	"eduhelper/edupage"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"strconv"
	"strings"
)

func Start() {
	gtk.Init(nil)
	var h handle

	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	if err != nil {
		log.Fatal("Unable to create window:", err)
	}

	h.Window = win

	win.SetTitle("Eduhelper")
	win.Connect("destroy", func() {
		h.quit()
	})

	err = h.loginPage()
	if err != nil {
		log.Fatal("Unable to create login page:", err)
	}

	win.SetDefaultSize(800, 600)
	win.ShowAll()

	gtk.Main()
}

func (h *handle) mainPage() error {
	grid, err := gtk.GridNew()
	if err != nil {
		return err
	}

	list, err := gtk.ListBoxNew()
	if err != nil {
		return err
	}

	err = list.SetProperty("activate-on-single-click", true)
	if err != nil {
		return err
	}

	ha, _ := gtk.AdjustmentNew(float64(1), float64(1), float64(1), float64(1), float64(1), float64(1))
	va, _ := gtk.AdjustmentNew(float64(1), float64(1), float64(1), float64(1), float64(1), float64(1))

	window, err := gtk.ScrolledWindowNew(ha, va)
	if err != nil {
		return err
	}

	window.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_EXTERNAL)
	window.Add(list)

	message, err := gtk.TextViewNew()
	if err != nil {
		return err
	}

	message.SetWrapMode(gtk.WRAP_CHAR) //FIX resize
	message.SetJustification(gtk.JUSTIFY_LEFT)
	message.SetEditable(false)
	message.SetCanFocus(false)
	h.message = message

	info, err := gtk.TextViewNew()
	if err != nil {
		return err
	}

	info.SetWrapMode(gtk.WRAP_CHAR) //FIX resize
	info.SetJustification(gtk.JUSTIFY_LEFT)
	info.SetEditable(false)
	info.SetCanFocus(false)
	h.info = info

	ntb, err := gtk.NotebookNew()
	if err != nil {
		return err
	}

	messageLabel, err := gtk.LabelNew("Message")
	if err != nil {
		return err
	}

	ntb.AppendPage(message, messageLabel)

	infoLabel, err := gtk.LabelNew("Info")
	if err != nil {
		return err
	}

	ntb.AppendPage(info, infoLabel)

	h.notebook = ntb

	grid.Attach(window, 0, 0, 1, 1)
	grid.Attach(ntb, 1, 0, 3, 1)

	grid.SetColumnHomogeneous(true)
	grid.SetRowHomogeneous(true)
	grid.SetColumnSpacing(5)

	tm, err := h.ehandle.GetTimeline()
	if err != nil {
		return err
	}

	items := tm.SortedTimelineItems(func(item edupage.TimelineItem) bool {
		return item.Type == edupage.TimelineMessage
	})

	h.listRows = items
	list.Connect("row-selected", h.rowSelect)

	for i, item := range items {
		var preview string
		if len(item.Text) > 20 {
			preview = string([]rune(item.Text)[0:20])
			i := strings.Index(preview, "\n")
			if i != -1 {
				preview = preview[0:i]
			}
			preview += "..."
		} else {
			preview = item.Text
		}

		rowBox, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 1)
		if err != nil {
			return err
		}

		title, err := gtk.LabelNew(item.OwnerName + ":")
		if err != nil {
			return err
		}

		rowBox.PackStart(title, true, true, 0)

		text, err := gtk.LabelNew(preview)
		if err != nil {
			return err
		}

		text.SetHAlign(gtk.ALIGN_START)
		rowBox.PackStart(text, true, true, 0)
		sep, err := gtk.SeparatorNew(gtk.ORIENTATION_HORIZONTAL)
		if err != nil {
			return err
		}

		rowBox.PackStart(sep, true, true, 0)

		row, err := gtk.ListBoxRowNew()
		if err != nil {
			return err
		}

		row.Add(rowBox)
		row.SetName(strconv.Itoa(i))

		list.Insert(row, i)
	}

	h.Window.Add(grid)
	h.Window.ShowAll()
	return nil
}

func (h *handle) loginPage() error {
	grid, err := gtk.GridNew()
	if err != nil {
		return err
	}

	label, err := gtk.LabelNew("")
	if err != nil {
		return err
	}

	label.SetLineWrap(true)
	label.SetMaxWidthChars(25)

	usernameE, err := gtk.EntryNew()
	if err != nil {
		return err
	}
	err = applyStyle(usernameE.ToWidget(), "flat")
	if err != nil {
		return err
	}

	passwordE, err := gtk.EntryNew()
	if err != nil {
		return err
	}
	passwordE.SetVisibility(false)

	err = applyStyle(passwordE.ToWidget(), "flat")
	if err != nil {
		return err
	}

	serverE, err := gtk.EntryNew()
	if err != nil {
		return err
	}

	err = applyStyle(serverE.ToWidget(), "flat")
	if err != nil {
		return err
	}

	serverL, err := gtk.LabelNew("Server")
	if err != nil {
		return err
	}

	usernameL, err := gtk.LabelNew("Username")
	if err != nil {
		return err
	}

	passwordL, err := gtk.LabelNew("Password")
	if err != nil {
		return err
	}

	loginB, err := gtk.ButtonNewWithLabel("Login")
	if err != nil {
		return err
	}
	loginB.Connect("button-press-event", func(*gtk.Button, *gdk.Event) {
		username, _ := usernameE.GetText()
		password, _ := passwordE.GetText()
		server, _ := serverE.GetText()

		e, err := edupage.Login(server, username, password)
		if err != nil {
			if err == edupage.AuthorizationError {
				label.SetText("Invalid password")
				passwordE.SetText("")
			} else {
				label.SetText("Error: " + err.Error())
			}
		} else {
			label.SetText("Success")
			h.ehandle = &e
			h.Window.Remove(grid)
			grid.Destroy()
			err = h.mainPage()
			if err != nil {
				log.Fatal("Unable to create main page after login:", err)
			}
		}

	})

	quitB, err := gtk.ButtonNewWithLabel("Quit")
	quitB.Connect("button-press-event", func(*gtk.Button, *gdk.Event) {
		h.quit()
	})
	if err != nil {
		return err
	}

	grid.SetRowSpacing(5)
	grid.SetColumnSpacing(10)
	grid.SetHAlign(gtk.ALIGN_CENTER)
	grid.SetVAlign(gtk.ALIGN_CENTER)

	grid.Attach(label, 0, 0, 3, 2)
	grid.Attach(usernameL, 0, 2, 1, 1)
	grid.Attach(usernameE, 1, 2, 1, 1)

	grid.Attach(passwordL, 0, 3, 1, 1)
	grid.Attach(passwordE, 1, 3, 1, 1)

	grid.Attach(serverL, 0, 4, 1, 1)
	grid.Attach(serverE, 1, 4, 1, 1)

	subgrid, err := gtk.GridNew()
	subgrid.Attach(quitB, 0, 0, 1, 1)
	subgrid.Attach(loginB, 1, 0, 1, 1)
	subgrid.SetColumnSpacing(25)
	subgrid.SetColumnHomogeneous(true)

	grid.Attach(subgrid, 0, 5, 2, 1)

	h.Window.Add(grid)

	return nil
}

func applyStyle(e *gtk.Widget, class string) error {
	s, err := e.GetStyleContext()
	s.AddClass(class)
	e.ResetStyle()
	return err
}
