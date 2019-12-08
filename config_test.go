package vconfig

import "testing"

func Test_Debug(t *testing.T) {
	SetDebug(true)
	if !Debug() {
		t.Errorf("Debut not set and should be.")
	}
}
