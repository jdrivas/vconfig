package vconfig

import "github.com/spf13/viper"

// Debug returns whether debug mode is set.
func Debug() bool {
	return viper.GetBool(DebugFlagKey)
}

// SetDebug allows you to turn on or off the debug mode.
func SetDebug(b bool) {
	viper.Set(DebugFlagKey, b)
}

// ToggleDebug toggles the flag and returns the new value.
func ToggleDebug() bool {
	viper.Set(DebugFlagKey, viper.GetBool(DebugFlagKey))
	return Debug()
}

// Verbose returs whether verbose mode is set.
func Verbose() bool {
	return viper.GetBool(VerboseFlagKey)
}

// ToggleVerbose toggles the flag and returns the new value.
func ToggleVerbose() bool {
	viper.Set(VerboseFlagKey, viper.GetBool(VerboseFlagKey))
	return Verbose()
}
