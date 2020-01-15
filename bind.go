package vconfig

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// BindFlag - structure for dealing with values that propogate from the
// app command line flags.
type BindFlag struct {
	Flag    *pflag.Flag
	BindKey string
	value   interface{}
}
type bindMap map[string]*BindFlag

var bfm = make(bindMap) // Keyed by flag
var bbm = make(bindMap) // keyed by binding

// Bind a flag key to a flag variable.
// The reason we're doing this is because the viper
// pFlags integration doesn't doesn't comprehend multiple
// invocations. So, the first a key is bound to a flag,
// a new BindEntry is created. After that if a bind entry already
// exists, only the flag will be updated.
func Bind(bk string, f *pflag.Flag) (bf *BindFlag) {
	if Debug() {
		pef()
		defer pxf()
	}

	// Always do both maps.
	// Create a new key if both have no values in the maps.
	// Else update the key that was found.
	var ok bool
	if bf, ok = bfm[f.Name]; !ok {
		if bf, ok = bbm[bk]; !ok { // Neither have been set
			if Debug() {
				fmt.Printf("New binding: %#v\n", bk)
			}
			bf = new(BindFlag)
		} else { // Out of Sync: bfm not set, bbm set.
			if Debug() {
				fmt.Printf("BindMaps are out of sync: keyed by flag missing, keyed by bind found: %#v\n",
					bf)
			}
		}
	} else if _, ok := bbm[bk]; !ok { // Out of Sync: bfm set, bbm not set.
		if Debug() {
			fmt.Printf("BindMaps are out of sync: keyed by flag found, keyed by binding missing: %#v\n",
				bf)
		}
	}

	// Always do BOTH maps.
	bf.BindKey = bk
	bf.Flag = f
	bfm[f.Name] = bf
	bbm[bk] = bf

	return bf
}

// GetBindFlags returns all the BindFlags registered.
func GetBindFlags() (bfs []*BindFlag) {
	for _, v := range bfm {
		bfs = append(bfs, v)
	}
	return bfs
}

// Set will set the viper variable and keep the
// value for later application during Apply.
func Set(bk string, value interface{}) {
	if bf, ok := bbm[bk]; ok {
		bf.value = value
	}
	viper.Set(bk, value)
}

// UpdateChangedFlags will look at each binding
// and if the associated flag has changed, udpate the bind value.
// This is intended to be used to capture values to bind to viper
// right after a parse of flags has acurred. You might then
// immediately call Apply() to cause the viper variables to take this new value.
// This is different behavior than ApplyFromFlags.
func UpdateChangedFlags() {
	if Debug() {
		pef()
		defer pxf()
	}
	for _, bf := range GetBindFlags() {
		if bf.Flag.Changed {
			if Debug() {
				fmt.Printf("Flag changed %q, setting bind value to: %q\n", bf.Flag.Name, bf.Flag.Value.String())
			}
			bf.setValueFrom(bf.Flag)
		}
	}
}

// Apply will set the viper variable with BindKey to the Value if
// there is a Value.
func Apply() {
	if Debug() {
		pef()
		defer pxf()
	}

	for _, bf := range bbm {
		if bf.value != nil {
			if Debug() {
				fmt.Printf("Setting viper value with key %#v with value %#v\n",
					bf.BindKey, bf.value)
			}
			viper.Set(bf.BindKey, bf.value)
		}
	}
}

// ApplyFromFlags will look at each flag in the flag set
// and if the flag changed and there is a Binding it will
// set the viper value to the value of the flag (not the BindValue).
// This is where precedence is maintained essentially allowing for
// a switch having flags take short-term preccedence over sets.
func ApplyFromFlags(pflags *pflag.FlagSet) {
	if Debug() {
		pef()
		defer pxf()
	}
	pflags.VisitAll(func(pf *pflag.Flag) {
		if Debug() {
			fmt.Printf("Visiting flag: %#v\n%s", pf.Name, flagString(pf))
		}
		if bf := bfm[pf.Name]; bf != nil { // if bound
			var v interface{}
			if pf.Changed { // and flag changed
				// Set the viper variable to the flag value.
				v = flagValue(pf)
				// Since we're going to set a viper value from a flag
				// we need to make sure there is a bind value to
				// replace it with at a later time through Apply.
				// So, check to see what the current value is,
				// and if there is no, use the default.
				// TODO: THIS BREAKS ANYMORE DYNAMIC UPDATES.
				// So if you're watching a config file and/or you're
				// using an environment variable those won't effect subsequent
				// viper queires of this value.
				if bf.value == nil { // If there is a current value use that.
					if viper.IsSet(bf.BindKey) {
						bf.value = viper.Get(bf.BindKey)
					} else {
						fdv := flagDefValue(pf) // otherwise use the default.
						bf.value = fdv
					}
				}
			} else if bf.value != nil { // or not changed and we have a bind value
				v = bf.value
			} // we don't care about the case where we're not changing by a flag and there is no bind value.
			// If we've set a viper value give it viper.
			if v != nil {
				if Debug() {
					fmt.Printf("Setting viper value %#v to %#v\n", bf.BindKey, v)
				}
				viper.Set(bf.BindKey, v)
			}
		}
	})
}

// ResetBindings will erase existing bindings.
// This is really used for Testing.
func ResetBindings() {
	bfm = make(bindMap)
	bbm = make(bindMap)
}

// GetBindFlagFor return BindFlag for the flag key.
func getBindFlagFor(fk string) *BindFlag {
	return bfm[fk]
}

// SetValueFrom sets the VindFLag value from a pfFlag.
// a flag.
func (bf *BindFlag) setValueFrom(f *pflag.Flag) {
	bf.value = flagValue(f)
}

// This is gratuitous and only used in test.
func flagForFlagKey(fk string) (bf *BindFlag, ok bool) {
	bf, ok = bfm[fk]
	return bf, ok
}

func flagValue(f *pflag.Flag) interface{} {
	return stringValue(f.Value.String(), f.Value.Type())
}

func flagDefValue(f *pflag.Flag) interface{} {
	t := f.Value.Type()
	return stringValue(f.DefValue, t)
}

// TODO: Need to do one for int and float.
func stringValue(v, t string) interface{} {
	switch t {
	case "string":
		return v
	case "bool":
		if strings.ToLower(v) == "true" {
			return true
		}
		return false
	default:
		return v
	}

}
