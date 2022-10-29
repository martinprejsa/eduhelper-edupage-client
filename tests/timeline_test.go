package tests

import (
	"eduhelper/edupage"
	"errors"
	"os"
	"testing"
)

func TestTimeline(t *testing.T) {
	h, err := edupage.Login(os.Getenv("EDUPAGE_SERVER"), os.Getenv("EDUPAGE_USERNAME"), os.Getenv("EDUPAGE_PASSWORD"))
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
}
