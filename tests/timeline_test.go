package tests

import (
	"eduhelper/edupage"
	"errors"
	"testing"
)

func TestTimeline(t *testing.T) {
	h, err := edupage.Login("SERVER", "USERNAME", "PASSWORD")
	if err != nil {
		t.Fatal(err)
	}

	ti, err := h.GetTimeline()
	if err != nil {
		t.Fatal(err)
	}

	if len(ti.TimelineItems) == 0 {
		t.Fatal(errors.New("no timeline items"))
	}

	for _, hw := range ti.Homeworks {
		if hw.TestID == "429860" {
			_, err := h.GetHomeworkAttachments(&hw)
			if err != nil && err != edupage.UnobtainableAttachments {
				t.Fatal(err)
			}
		}
	}

}
