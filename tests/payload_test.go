package tests

import (
	"eduhelper/edupage"
	"testing"
)

func TestPayload(t *testing.T) {
	p, err := edupage.CreatePayload(map[string]string{
		"test": "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	if p.Encode() != "eqacs=q75pYGlWjPO7-0oFkF1_AjUfGQ8%3D&eqap=dGVzdD10ZXN0&eqaz=1" {
		t.Fatal("control payload mismatch")
	}
}
