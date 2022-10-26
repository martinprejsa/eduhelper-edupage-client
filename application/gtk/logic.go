package gtk

import (
	"eduhelper/edupage"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pkg/browser"
	"path"
	"strconv"
)

type handle struct {
	Window   *gtk.Window
	ehandle  *edupage.Handle
	notebook *gtk.Notebook
	listRows []edupage.TimelineItem
	message  *gtk.TextView
	info     *gtk.TextView
}

func (h *handle) quit() {
	h.ehandle = nil
	h.notebook = nil
	h.message = nil
	h.message = nil
	h.info = nil
	gtk.MainQuit()
}

func (h *handle) rowSelect(_ *gtk.ListBox, row *gtk.ListBoxRow) {
	h.notebook.RemovePage(2)
	h.notebook.SetCurrentPage(0)
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

	text := h.listRows[index].Text
	if h.listRows[index].Data.Value["receipt"] == "1" {
		text += "\n\n"
		text += h.listRows[index].Data.Value["messageContent"].(string)
	}
	messageBuffer.SetText(text)
	h.message.SetBuffer(messageBuffer)

	itb, err := gtk.TextTagTableNew()
	if err != nil {
		return
	}

	infoBuffer, err := gtk.TextBufferNew(itb)
	if err != nil {
		return
	}

	infoBuffer.SetText(
		"Author: " + h.listRows[index].OwnerName + "\n" +
			"Sent to: " + h.listRows[index].UserName + "\n" +
			"Created at: " + h.listRows[index].TimeAdded.Format(edupage.TimeFormat) + "\n",
	)
	h.info.SetBuffer(infoBuffer)

	attachments := h.listRows[index].GetAttachments()

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

		h.notebook.AppendPage(attachmentList, attachmentsLabel)

		for key, val := range attachments {
			btn, err := gtk.ButtonNewWithLabel(key)
			if err != nil {
				return
			}

			btn.Connect("button-press-event", func(*gtk.Button, *gdk.Event) {
				_ = browser.OpenURL("https://" + path.Join(edupage.Server, val)) //TODO LOG
			})
			row, err := gtk.ListBoxRowNew()
			if err != nil {
				return
			}

			row.Add(btn)
			attachmentList.Add(row)
			h.Window.ShowAll()
			index++
		}
	}
}
