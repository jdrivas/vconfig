package vconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/juju/ansiterm"
	"github.com/spf13/viper"
)

// Debug returns whether debug mode is set.
func Debug() bool {
	return viper.GetBool(DebugKey)
}

// SetDebug allows you to turn on or off the debug mode.
func SetDebug(b bool) {
	viper.Set(DebugKey, b)
}

// ToggleDebug toggles the flag and returns the new value.
func ToggleDebug() bool {
	viper.Set(DebugKey, !viper.GetBool(DebugKey))
	return Debug()
}

// Verbose returs whether verbose mode is set.
func Verbose() bool {
	return viper.GetBool(VerboseKey)
}

// ToggleVerbose toggles the flag and returns the new value.
func ToggleVerbose() bool {
	viper.Set(VerboseKey, !viper.GetBool(VerboseKey))
	return Verbose()
}

// pef is an entry tracing printout to stdout. It will print the word Enter, a funciton name, file name and line number.
// It is designed to be used as the first line of a function, perhaps bracketed by a debug check.
func pef() {
	fc, fl, ln := locString(2)
	w := ansiterm.NewTabWriter(os.Stdout, 6, 2, 1, ' ', 0)
	fmt.Fprintf(w, "Enter\t%s\t%s:%d\n", fc, fl, ln)
	w.Flush()
}

// pxf is an exit tracing pritout to stdout. It will print out the word exit, a function name, file anda line number.
// It is designed to be used just after a call to Pef() with a defer.
// Pef()
// defer Pxf()
func pxf() {
	fc, fl, ln := locString(2)
	w := ansiterm.NewTabWriter(os.Stdout, 6, 2, 1, ' ', 0)
	fmt.Fprintf(w, "Exit\t%s\t%s:%d\n", fc, fl, ln)
	w.Flush()
}

func locString(d int) (fnc, file string, line int) {
	if pc, fl, l, ok := runtime.Caller(d); ok {
		f := runtime.FuncForPC(pc)
		fnc = filepath.Base(f.Name())
		file = filepath.Base(fl)
		line = l
	}
	return fnc, file, line
}
