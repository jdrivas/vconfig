package vconfig

/*
* Keys to look up values in Viper configuration.
 */

// TODO: Formalize the yaml file structure, or at least document it here.

// YAML Variables which show up in viper, but managed here.

/*
 These are top level variables in, say, a YAML file.
 e.g. debug = true
 not:
 vconfig:
		 debug: true
*/
const (
	DebugKey   = "debug"   // bool
	VerboseKey = "verbose" // bool
)
