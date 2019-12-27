package vconfig

import (
	"fmt"
	"testing"
)

func Test_Debug(t *testing.T) {
	SetDebug(true)
	if !Debug() {
		t.Errorf("Debug not set and should be.")
	}

	ToggleDebug()
	if Debug() {
		t.Error("Debug is set and should have toggled off.")
	}
}

func Test_AppName(t *testing.T) {
	testAppName := "TestApp"
	AppName = testAppName
	InitConfig()

	if ConfigFileRoot != testAppName {
		t.Errorf("Bad config file root. Got: %s, Expected %s.", ConfigFileRoot, testAppName)
	}

	hFileName := fmt.Sprintf(".%s_history", testAppName)
	if HistoryFile != hFileName {
		t.Errorf("Bad history file name. Got: %s, Expected %s.", HistoryFile, hFileName)

	}
}
