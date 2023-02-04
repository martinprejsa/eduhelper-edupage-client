package gtk

import (
	"eduhelper/edupage"
	"fmt"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/browser"
	"path"
	"strconv"
)

type handle struct {
	Window   *gtk.Window
	Explorer *gtk.ScrolledWindow
	Viewer   *gtk.ScrolledWindow
	ehandle  *edupage.Handle
	listRows []edupage.TimelineItem
}

func (h *handle) rowSelect(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	if child, err := h.Viewer.GetChild(); child != nil && err == nil {
		child.ToWidget().Destroy()
	}

	message, err := gtk.TextViewNew()
	if err != nil {
		return
	}

	message.SetWrapMode(gtk.WRAP_CHAR)
	message.SetJustification(gtk.JUSTIFY_LEFT)
	message.SetEditable(false)
	message.SetCanFocus(false)

	info, err := gtk.TextViewNew()
	if err != nil {
		return
	}

	info.SetWrapMode(gtk.WRAP_CHAR)
	info.SetJustification(gtk.JUSTIFY_LEFT)
	info.SetEditable(false)
	info.SetCanFocus(false)

	ntb, err := gtk.NotebookNew()
	if err != nil {
		return
	}

	messageLabel, err := gtk.LabelNew("Message")
	if err != nil {
		return
	}

	ntb.AppendPage(message, messageLabel)

	infoLabel, err := gtk.LabelNew("Info")
	if err != nil {
		return
	}

	ntb.AppendPage(info, infoLabel)
	h.Viewer.Add(ntb)

	/** FILLING **/

	name, err := row.GetName()
	if err != nil {
		return
	}

	index, err := strconv.Atoi(name)
	if err != nil {
		return
	}

	mtb, err := gtk.TextTagTableNew()
	if err != nil {
		return
	}

	messageBuffer, err := gtk.TextBufferNew(mtb)
	if err != nil {
		return
	}

	itb, err := gtk.TextTagTableNew()
	if err != nil {
		return
	}

	infoBuffer, err := gtk.TextBufferNew(itb)
	if err != nil {
		return
	}

	ctx := h.listRows[index]

	infoBuffer.SetText(
		"Author: " + ctx.OwnerName + "\n" +
			"Sent to: " + ctx.UserName + "\n" +
			"Created at: " + ctx.TimeAdded.Format(edupage.TimeFormat) + "\n",
	)
	info.SetBuffer(infoBuffer)

	var attachments map[string]string
	if ctx.Type == edupage.TimelineMessage {
		text := ctx.Text
		if ctx.Data.Value["receipt"] == "1" {
			text += "\n\n"
			text += ctx.Data.Value["messageContent"].(string)
		}

		messageBuffer.SetText(text)

		attachments, err = ctx.GetAttachments()
	} else if ctx.Type == edupage.TimelineHomework {
		hw, err := ctx.ToHomework()
		if err != nil {
			fmt.Println(err)
			return
		}
		attachments, err = h.ehandle.FetchHomeworkAttachments(&hw)
		if err != nil {
			fmt.Println(err)
			return
		}
		text := ctx.Data.Value["nazov"].(string)
		messageBuffer.SetText(text)
	}

	if len(attachments) != 0 {
		attachmentsLabel, err := gtk.LabelNew("Attachments")
		if err != nil {
			return
		}

		attachmentList, err := gtk.ListBoxNew()
		if err != nil {
			return
		}

		attachmentList.SetActivateOnSingleClick(false)
		ntb.AppendPage(attachmentList, attachmentsLabel)
		for key, val := range attachments {
			btn, err := gtk.ButtonNewWithLabel(key)
			if err != nil {
				return
			}

			btn.Connect("button-press-event", func(*gtk.Button, *gdk.Event) {
				_ = browser.OpenURL("https://" + path.Join(edupage.Server, val)) //TODO LOG
			})
			row, err := gtk.ListBoxRowNew()
			row.SetActivatable(false)
			row.SetSelectable(false)
			if err != nil {
				return
			}
			row.Add(btn)
			attachmentList.Add(row)

			index++
		}
	}

	message.SetBuffer(messageBuffer)
	h.Window.ShowAll()
	return
}
