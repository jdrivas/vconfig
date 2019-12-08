package vconfig

import "github.com/spf13/viper"

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
	viper.Set(DebugKey, viper.GetBool(DebugKey))
	return Debug()
}

// Verbose returs whether verbose mode is set.
func Verbose() bool {
	return viper.GetBool(VerboseKey)
}

// ToggleVerbose toggles the flag and returns the new value.
func ToggleVerbose() bool {
	viper.Set(VerboseKey, viper.GetBool(VerboseKey))
	return Verbose()
}
