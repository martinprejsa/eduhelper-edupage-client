package gtk

import (
	"eduhelper/edupage"
	"eduhelper/utils"
	"errors"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"log"
	"os"
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
		gtk.MainQuit()
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

	eha, _ := gtk.AdjustmentNew(float64(1), float64(1), float64(1), float64(1), float64(1), float64(1))
	eva, _ := gtk.AdjustmentNew(float64(1), float64(1), float64(1), float64(1), float64(1), float64(1))

	explorer, err := gtk.ScrolledWindowNew(eha, eva)
	if err != nil {
		return err
	}

	explorer.SetPolicy(gtk.POLICY_NEVER, gtk.POLICY_EXTERNAL)
	explorer.Add(list)

	vha, _ := gtk.AdjustmentNew(float64(1), float64(1), float64(1), float64(1), float64(1), float64(1))
	vva, _ := gtk.AdjustmentNew(float64(1), float64(1), float64(1), float64(1), float64(1), float64(1))

	viewer, err := gtk.ScrolledWindowNew(vha, vva)
	if err != nil {
		return err
	}

	viewer.SetPolicy(gtk.POLICY_AUTOMATIC, gtk.POLICY_AUTOMATIC)

	grid.Attach(explorer, 0, 0, 1, 1)
	h.Explorer = explorer
	grid.Attach(viewer, 1, 0, 3, 1)
	h.Viewer = viewer

	grid.SetColumnHomogeneous(true)
	grid.SetRowHomogeneous(true)
	grid.SetColumnSpacing(5)

	tm, err := h.ehandle.GetTimeline()
	if err != nil {
		return err
	}

	items := tm.SortedTimelineItems(func(item edupage.TimelineItem) bool {
		return item.Type == edupage.TimelineMessage ||
			(item.Type == edupage.TimelineHomework && item.IsHomeworkWithAttachments())
	})

	h.listRows = items
	list.Connect("row-selected", h.rowSelect)

	for i, item := range items {
		var preview string
		if item.Type == edupage.TimelineHomework {
			preview = item.Data.Value["nazov"].(string)
		}

		if item.Type == edupage.TimelineMessage {
			preview = item.Text
		}

		if len(preview) > 20 {
			preview = preview[0:20]
			i := strings.Index(preview, "\n")
			if i != -1 {
				preview = preview[0:i]
			}
			preview += "..."
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
	_ = os.MkdirAll(utils.GetRootDir(), os.FileMode(0700)) //TODO: log
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

	quitB, err := gtk.ButtonNewWithLabel("Quit")
	quitB.Connect("button-press-event", func(*gtk.Button, *gdk.Event) {
		gtk.MainQuit()
	})
	if err != nil {
		return err
	}

	rememberCheckbox, err := gtk.CheckButtonNewWithLabel("Remember credentials (unsafe)")
	if err != nil {
		return err
	}

	saved := false
	server, username, password, err := loadCredentials()
	if err == nil {
		serverE.SetText(server)
		usernameE.SetText(username)
		passwordE.SetText(password)
		rememberCheckbox.SetActive(true)
		saved = true
	}

	h.Window.Connect("key-pressed", func(g *gtk.Window, event *gdk.EventKey) {
		if event.KeyVal() == gdk.KEY_Escape {
			quitB.Clicked()
		} else if event.KeyVal() == gdk.KEY_ISO_Enter {
			loginB.Clicked()
		}
	})

	loginB.Connect("button-press-event", func(*gtk.Button, *gdk.Event) {
		username, _ := usernameE.GetText()
		password, _ := passwordE.GetText()
		server, _ := serverE.GetText()

		e, err := edupage.Login(server, username, password)
		if err != nil {
			if err == edupage.AuthorizationError {
				label.SetText("Invalid credentials")
				passwordE.SetText("")
			} else {
				label.SetText("Error: " + err.Error())
			}
		} else {
			if !saved && rememberCheckbox.GetActive() {
				_ = saveCredentials(server, username, password) //TODO: log
			}

			if !rememberCheckbox.GetActive() {
				_ = os.Remove(utils.GetCredentialsFilePath()) //TODO: log
			}

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

	grid.Attach(rememberCheckbox, 0, 5, 2, 1)
	grid.Attach(subgrid, 0, 6, 2, 1)

	h.Window.Add(grid)

	return nil
}

// TODO write test
func saveCredentials(server, username, password string) error {
	var escape = func(input string) string {
		return strings.ReplaceAll(input, ":", "\\:")
	}
	str := escape(server) + ":" + escape(username) + ":" + escape(password)
	err := os.WriteFile(utils.GetCredentialsFilePath(), []byte(str), os.FileMode(0700))
	if err != nil {
		return err
	}
	return nil
}

// TODO write test
// server, username, password
func loadCredentials() (string, string, string, error) {
	data, err := os.ReadFile(utils.GetCredentialsFilePath())
	if err != nil {
		return "", "", "", err
	}
	if len(data) == 0 {
		return "", "", "", errors.New("no credentials present")
	}

	str := string(data)
	items := strings.Split(str, ":")
	for _, item := range items {
		item = strings.ReplaceAll(item, "\\:", ":")
	}
	if len(items) != 3 {
		return "", "", "", errors.New("invalid credentials")
	}

	return items[0], items[1], items[2], nil
}

func applyStyle(e *gtk.Widget, class string) error {
	s, err := e.GetStyleContext()
	s.AddClass(class)
	e.ResetStyle()
	return err
}
