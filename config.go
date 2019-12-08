package vconfig

import (
	"fmt"
	"os"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	AppName = "application"
)

var (
	ConfigFileName string
	ConfigFileRoot = fmt.Sprintf("%s", AppName)
	HistoryFile    = fmt.Sprintf(".%s_history", AppName)
)

/*
* Keys to look up values in Viper configuration.
 */

// TODO: Formalize the yaml file structure, or at least document it here.

// YAML Variables which show up in viper, but managed here.
const (
	DebugKey         = "debug"         // bool
	VerboseKey       = "verbose"       // bool
	ScreenProfileKey = "screenProfile" // string
	ScreenDarkValue  = "dark"
)

// Flags These are the long form flag values for command line flags.
const (
	ConfigFlagKey  = "config"
	VerboseFlagKey = "verbose"
	DebugFlagKey   = "debug"
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {

	// fmt.Printf("%s\n", t.Title("InitConfig")) // Can't bracket with util.Debug as Debug uses config.

	// Fin a config file
	if ConfigFileName != "" {
		viper.SetConfigFile(ConfigFileName)
	} else {
		viper.SetConfigName(ConfigFileRoot)

		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra_test" (without extension).
		viper.AddConfigPath(".")
		viper.AddConfigPath(home)
	}

	viper.AutomaticEnv() // read in environment variables that match

	// Read in the config file.
	if err := viper.ReadInConfig(); err == nil {
		if Debug() {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	} else {
		fmt.Printf("Error loading config file: %s - %v\n", viper.ConfigFileUsed(), err)
	}

	// fmt.Printf("%s\n", t.Title("InitConfig - exit"))

}
