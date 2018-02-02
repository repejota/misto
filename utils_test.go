package misto

import (
	"testing"
)

func TestUtilsDummy(t *testing.T) {
}

func TestStripCtlAndExtFromUnicode(t *testing.T) {
	orig := "foo bar"
	new := StripCtlAndExtFromUnicode(orig)
	if orig != new {
		t.Fatalf("Expecetd %q but got %q", orig, new)
	}
}
