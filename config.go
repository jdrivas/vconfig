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
	ConfigFileRoot string
	HistoryFile    string
)

// InitConfig reads in config file and ENV variables if set.
func InitConfig() {

	if Debug() {
		pef() 
		defer pxf()
	}

	if ConfigFileRoot == "" {
		ConfigFileRoot = fmt.Sprintf("%s", AppName)
	}

	if HistoryFile == "" {
		HistoryFile = fmt.Sprintf(".%s_history", AppName)
	}

	// Find a config file
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

}
