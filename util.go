package vconfig

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/juju/ansiterm"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Debug returns whether debug mode is set.
func Debug() bool {
	return viper.GetBool(DebugKey)
}

// SetDebug allows you to turn on or off the debug mode.
func SetDebug(b bool) {
	Set(DebugKey, b)
}

// ToggleDebug toggles the flag and returns the new value.
func ToggleDebug() bool {
	Set(DebugKey, !viper.GetBool(DebugKey))
	return Debug()
}

// Verbose returs whether verbose mode is set.
func Verbose() bool {
	return viper.GetBool(VerboseKey)
}

// ToggleVerbose toggles the flag and returns the new value.
func ToggleVerbose() bool {
	Set(VerboseKey, !viper.GetBool(VerboseKey))
	return Verbose()
}

func flagString(pf *pflag.Flag) string {
	var b strings.Builder
	w := ansiterm.NewTabWriter(&b, 4, 4, 2, ' ', 0)
	fmt.Fprintf(w, flagHeader()+"\n")
	fmt.Fprintf(w, flagEntry(pf)+"\n")
	w.Flush()
	return b.String()
}

func flagHeader() string {
	return "Name\tShort\tValue\tType\tDefValue\tChanged"
}

func flagEntry(f *pflag.Flag) string {
	if f != nil {
		return fmt.Sprintf("%s\t%s\t%s\t%s\t%s\t%t",
			f.Name, f.Shorthand, f.Value.String(), f.Value.Type(), f.DefValue, f.Changed)
	}
	return "<No Flag>\t-\t-\t-\t-\t-"
}

// pef is an entry tracing printout to stdout. It will print the word Enter, a funciton name, file name and line number.
// It is designed to be used as the first line of a function, perhaps bracketed by a debug check.
func pef() {
	fc, fl, ln := loc(1)
	w := ansiterm.NewTabWriter(os.Stdout, 6, 2, 1, ' ', 0)
	fmt.Fprintf(w, "Enter\t%s\t%s:%d\n", fc, fl, ln)
	w.Flush()
}

// pxf is an exit tracing pritout to stdout. It will print out the word exit, a function name, file anda line number.
// It is designed to be used just after a call to Pef() with a defer.
// Pef()
// defer Pxf()
func pxf() {
	fc, fl, ln := loc(1)
	w := ansiterm.NewTabWriter(os.Stdout, 6, 2, 1, ' ', 0)
	fmt.Fprintf(w, "Exit\t%s\t%s:%d\n", fc, fl, ln)
	w.Flush()
}

func locString(d int) string {
	fc, fl, ln := loc(d + 1)
	return fmt.Sprintf("%s()  %s:%d", fc, fl, ln)
}

func loc(d int) (fnc, file string, line int) {
	if pc, fl, l, ok := runtime.Caller(d + 1); ok {
		f := runtime.FuncForPC(pc)
		fnc = filepath.Base(f.Name())
		file = filepath.Base(fl)
		line = l
	}
	return fnc, file, line
}
