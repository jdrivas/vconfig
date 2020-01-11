package vconfig

import (
	"testing"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type flagExpected struct {
	fKey, fV string
}

type set struct {
	key, value string
}

type flag struct {
	bk, fk, fks, fv, fd, fh string
	fVar                    *string
	e                       flagExpected
}

type bflag struct {
	bk, fk, fks, fh, fv string
	fd                  bool
	fVar                *bool
}

var f1 = flag{
	bk:  "filename",                // binding key
	fk:  "filename",                // flag key
	fks: "f",                       // short flag key
	fv:  "somefile",                // flag value
	fd:  "",                        // flag default value
	fh:  "help message about file", // help message
	e: flagExpected{
		fKey: "filename",
		fV:   "somefile",
	},
}
var s1 = set{key: "filename", value: "someotherfile"}

var f2 = flag{
	bk:  "screen",                    // binding key
	fk:  "screenFlag",                // flag key
	fks: "s",                         // short flag key
	fv:  "screenName",                // flag value
	fd:  "",                          // flag default value
	fh:  "help message about screen", // help message
	e: flagExpected{
		fKey: "screenFlag",
		fV:   "screenName",
	},
}

var s2 = set{key: "screen", value: "darkScreen"}

// Create an argument list for parsing from an array of flags.
func argsFromFlags(app string, flags []flag) (args []string) {
	args = append(args, app)
	for _, f := range flags {
		args = append(args, "--"+f.fk, f.fv)
	}
	return args
}

// Register flags into a flagset, and then bind them
func registerAndBindFlags(flags []flag, pflags *pflag.FlagSet) {
	for _, f := range flags {
		f.fVar = pflags.String(f.fk, f.fd, f.fh)
		Bind(f.bk, pflags.Lookup(f.fk))
	}

}

// After parse, get the values that have changed from a flag set and
// set the bind value to flag value from the parse.
func getChangedValues(flags []flag, pflags *pflag.FlagSet, t *testing.T) {
	for _, f := range flags {
		if pf := pflags.Lookup(f.fk); pf != nil {
			if pf.Changed {
				if bf, ok := flagForFlagKey(pf.Name); ok {
					bf.value = flagValue(pf)
					t.Logf("Changed flag Value is %#v", bf.value)
				} else {
					t.Errorf("Failed to find bind flag, for flag %s", pf.Name)
				}
			}
		} else {
			t.Errorf("Couldn't find a pflag flag for expected flag: \"%s\"", f.fk)
		}
	}

}

// This is just to make sure the integration actually works.
func TestPflag(t *testing.T) {
	type tc struct {
		name  string
		flags []flag
		args  []string
	}

	// The actual test Cases
	cases := []tc{
		{name: "one flag", flags: []flag{f1}},
		{name: "two flags", flags: []flag{f1, f2}},
	}

	// Build out the test cases.
	for i, c := range cases {
		cases[i].args = argsFromFlags("app", c.flags)
	}

	for _, c := range cases {
		ResetBindings()
		t.Run(c.name, func(t *testing.T) {
			pflags := pflag.NewFlagSet(c.name, pflag.PanicOnError)
			registerAndBindFlags(c.flags, pflags)
			pflags.Parse(c.args[1:])

			for _, f := range c.flags {
				pf := pflags.Lookup(f.fk)
				if pf != nil {
					if pf.Value.String() != f.e.fV {
						t.Errorf("Didn't get the expected value parsed into the flag.")
					}
				} else {
					t.Errorf("Didn't find a flag that should have been added: \"%s\"", f.fk)
				}
			}

		})
	}
}

func TestParseWithBind(t *testing.T) {
	type tc struct {
		name   string
		flags  []flag
		mapLen int
		args   []string
	}

	// The actual test Cases
	cases := []tc{
		{name: "one flag", flags: []flag{f1}},
		{name: "two flags", flags: []flag{f2, f1}},
	}

	// Build out the test cases.
	for i, c := range cases {
		cases[i].mapLen = len(c.flags)
		cases[i].args = argsFromFlags("app", c.flags)
	}

	for _, c := range cases {
		ResetBindings() // reset the bind map.
		t.Run(c.name, func(t *testing.T) {

			// Setup
			//

			pflags := pflag.NewFlagSet(c.name, pflag.PanicOnError)
			registerAndBindFlags(c.flags, pflags)
			pflags.Parse(c.args[1:])

			// This is a good check if things are going wrong.
			for _, f := range c.flags {
				pf := pflags.Lookup(f.fk)
				if pf != nil {
					if pf.Value.String() != f.e.fV {
						t.Errorf("Didn't get the expected value parsed into the flag.")
					}
				} else {
					t.Errorf("Didn't find a flag that should have been added: %#v", f.fk)
				}
			}

			getChangedValues(c.flags, pflags, t)

			// Tests
			if len(bfm) != c.mapLen || len(bbm) != c.mapLen {
				t.Errorf("bindMap not the right size. Expected %d, got bfm = %d  and bbm = %d entries.",
					c.mapLen, len(bfm), len(bbm))
			}

			for _, f := range c.flags {

				t.Logf("Checking flag: %s", f.fk)

				// Make sure the two lookups (bfm, bbm), don't get out of sync.
				if _, ok := bbm[f.bk]; ok {
					if _, ok := flagForFlagKey(f.fk); !ok {
						t.Errorf("The bind map and the flag map are out of sync: Bind has key, flag doesn't.")
						t.Logf("Bind Map: %#v\n", bbm)
						t.Logf("Flag Map: %#v\n", bfm)
					}
				} else {
					if _, ok := flagForFlagKey(f.fk); ok {
						t.Errorf("The bind map and the flag map are out of sync: Bind doesn't have key, flag does.")
						t.Logf("Bind Map: %#v\n", bbm)
						t.Logf("Flag Map: %#v\n", bfm)
					} else {
						t.Errorf("Neither the bind nor the flag maps got the expected keys.")
						t.Logf("Bind Map: %#v\n", bbm)
						t.Logf("Flag Map: %#v\n", bfm)
					}
				}

				// Ensure that the bind values are the correct ones.
				e := f.e
				if bf, ok := flagForFlagKey(e.fKey); !ok {
					t.Errorf("Failed to get a BindFlag for flag key: %s", e.fKey)
				} else {
					if bf.value != nil {
						// Assume these are strings.
						v := bf.value.(string)
						if v != e.fV {
							t.Errorf("BindFlag and flag have different values. Got:%#v, Expected: %#v",
								v, e.fV)
						}
					} else {
						t.Errorf("Didn't get a value bound to BindFlag, expected: %#v", bf.Flag.Value.String())
					}
				}
			}
		})
	}
}

func TestBindSet(t *testing.T) {
	type tc struct {
		name  string
		flags []flag
		args  []string
		vSets []set
	}

	cases := []tc{
		{name: "one flag with var", flags: []flag{f1}, vSets: []set{s1}},
		{name: "two flags", flags: []flag{f2, f1}, vSets: []set{s1, s2}},
	}

	// Complete test cases.
	for i, c := range cases {
		cases[i].args = argsFromFlags("app", c.flags)
	}

	for _, c := range cases {
		ResetBindings()
		viper.Reset()
		pflags := pflag.NewFlagSet(c.name, pflag.PanicOnError)
		registerAndBindFlags(c.flags, pflags)
		pflags.Parse(c.args[1:])
		// Update bind values from the parse
		getChangedValues(c.flags, pflags, t)

		// Set some viper variables
		for _, s := range c.vSets {
			Set(s.key, s.value)
		}

		t.Run(c.name, func(t *testing.T) {
			for _, s := range c.vSets {
				// Check bind key.
				if bk, ok := bbm[s.key]; ok {
					bV := bk.value.(string)
					if bV != s.value {
						t.Errorf("Bound value mismatched after set. Got: %#v, expected %#v", bV, s.value)
					}
				} else {
					t.Errorf("Couldn't find bind key for expected bound value: %#v", s.key)
				}
				// Check viper
				vs := viper.GetString(s.key)
				if vs != s.value {
					t.Errorf("Viper didn't return the expected value. Got: %#v, Expected: %#v", vs, s.value)
				}
			}
		})

	}

}

func TestApply(t *testing.T) {
	type tc struct {
		name  string
		flags []flag
		args  []string
		vSets []set
	}

	cases := []tc{
		{name: "one flag with var", flags: []flag{f1}, vSets: []set{s1}},
		{name: "two flags", flags: []flag{f2, f1}, vSets: []set{s1, s2}},
	}

	// Complete test cases.
	for i, c := range cases {
		cases[i].args = argsFromFlags("app", c.flags)
	}

	for _, c := range cases {
		ResetBindings()
		viper.Reset()
		pflags := pflag.NewFlagSet(c.name, pflag.PanicOnError)
		registerAndBindFlags(c.flags, pflags)
		pflags.Parse(c.args[1:])
		// Update bind values from the parse
		getChangedValues(c.flags, pflags, t)

		t.Run(c.name, func(t *testing.T) {
			Apply()

			// Values should be set to the flag value.
			for _, f := range c.flags {
				v := viper.GetString(f.bk)
				if v != f.fv {
					t.Errorf("Incorrect value after apply. Expected: %#v, Got: %#v", f.fv, v)
				}
			}

			// Now set the values
			for _, s := range c.vSets {
				Set(s.key, s.value)
			}

			Apply()

			// Check to see if the apply kept the sets.
			for _, s := range c.vSets {
				v := viper.GetString(s.key)
				if v != s.value {
					t.Errorf("Incorrect value after set then apply. Expected: %#v, Got: %#v", s.value, v)
				}
			}
		})
	}
}

func TestApplySpecial(t *testing.T) {

	flags := []flag{f1, f2}
	set := s1

	ResetBindings()

	pflags := pflag.NewFlagSet("ApplySpecial", pflag.PanicOnError)
	registerAndBindFlags(flags, pflags)

	args := argsFromFlags("app", flags)
	pflags.Parse(args[1:])

	// Pick up the flags and use them as overides.
	getChangedValues(flags, pflags, t)

	// Except for variables that have been set:
	Set(set.key, set.value)

	// Then apply the state to variables
	Apply()

	// Let's make sure everythig is as we expect it.
	for _, f := range flags {

		if bf, ok := flagForFlagKey(f.fk); !ok {
			t.Errorf("Couldn't find expected BindFlag for flag key: %#v", f.fk)
		} else {
			// Viper Value
			vV := viper.GetString(bf.BindKey)
			// This was set, so we expect the set value
			if bf.BindKey == set.key {
				if vV != set.value {
					t.Errorf("Viper value is not the set value. Got: %#v, Expected: %#v", vV, set.value)
				}
			} else {
				// Not set, so this should be the flag value.
				if vV != f.fv {
					t.Errorf("Viper value is not the flag value. Got: %#v, Expected; %#v", vV, f.fv)
				}
			}
		}
	}
}

func TestGet(t *testing.T) {

	// Configure
	flags := []flag{f1, f2}

	// Setup
	ResetBindings()

	pflags := pflag.NewFlagSet("ApplySpecial", pflag.PanicOnError)
	registerAndBindFlags(flags, pflags)

	args := argsFromFlags("app", flags)
	pflags.Parse(args[1:])

	bfs := GetBindFlags()

	// Tests
	//

	// Got the right number.
	if len(bfs) != len(flags) {
		t.Errorf("Wrong number of BindFlags returned by GetBindFlags(). Got: %#v, Expedted: %#v",
			len(bfs), len(flags))
	}

	// All flags should have changed (and the attached Flag knows this.)
	for _, bf := range bfs {
		if !bf.Flag.Changed {
			t.Errorf("Flag in BindFlag part of parse and assigned a value, but not marked changed: %#v",
				bf.Flag.Name)
		}
	}

}

func TestBindDebug(t *testing.T) {
	type bflag struct {
		bk, fk, fks, fh, fv string
		fd                  bool
		fVar                *bool
	}

	f := bflag{
		bk:  DebugKey,                   // binding key
		fk:  DebugKey,                   // flag key
		fks: "d",                        // short flag key
		fv:  "",                         // flag value
		fd:  false,                      // flag default value
		fh:  "help message about debug", // help message
	}

	// Setup
	ResetBindings()

	pflags := pflag.NewFlagSet("ApplySpecial", pflag.PanicOnError)
	f.fVar = pflags.Bool(f.fk, f.fd, f.fh)
	Bind(f.bk, pflags.Lookup(f.fk))

	pflags.Parse([]string{"--" + f.fk})

	// Check for changed flag values and set the binding variable.
	bfs := GetBindFlags()
	for _, bf := range bfs {
		if bf.Flag.Changed {
			bf.SetValueFrom(bf.Flag)
		}
	}

	// Now Apply to Viper.
	Apply()

	if Debug() != true {
		t.Errorf("Debug should have been set by flag.")
	}

	ToggleDebug()

	if Debug() != false {
		t.Errorf("Debug should have been reset by flag.")
	}

}

func TestBindDebugFirstInteractive(t *testing.T) {

	f := bflag{
		bk:  DebugKey,                   // binding key
		fk:  DebugKey,                   // flag key
		fks: "d",                        // short flag key
		fv:  "",                         // flag value
		fd:  false,                      // flag default value
		fh:  "help message about debug", // help message
	}

	// Setup
	ResetBindings()

	pflags := pflag.NewFlagSet("ApplySpecial", pflag.PanicOnError)

	f.fVar = pflags.Bool(f.fk, f.fd, f.fh)
	Bind(f.bk, pflags.Lookup(f.fk))

	// Main command line.
	pflags.Parse([]string{"--" + f.fk})

	// Check for changed flag values and set the binding variable.
	bfs := GetBindFlags()
	for _, bf := range bfs {
		if bf.Flag.Changed {
			bf.SetValueFrom(bf.Flag)
		}
	}

	// Now Apply to Viper.
	Apply()

	// Test
	if Debug() != true {
		t.Errorf("Debug should have been set by flag.")
	}

	ToggleDebug()

	if Debug() != false {
		t.Errorf("Debug should have been reset by flag.")
	}

	// Interactive Command line: THE FLAG SHOULD NOT BE DURABLE.
	pflags.Parse([]string{"--" + f.fk})

	// Test
	//
	visit := false
	pflags.Visit(func(pf *pflag.Flag) {
		if pf.Name == f.fk {
			visit = true
			if bf := bfm[pf.Name]; bf != nil {
				viper.Set(bf.BindKey, flagValue(pf))
			}
		}
	})

	if visit == false {
		t.Errorf("Never visited the debug flag after second parse.")
	}

	if Debug() != true {
		t.Errorf("Debug false, expected true after second parse.")
	}

	Apply()

	if Debug() != false {
		t.Errorf("Debug true, should be calse after second parse and apply.")
	}

}

func TestMultipleBindsCmdLineFlag(t *testing.T) {
	// SetDebug(true)
	// Configure
	flags := []flag{f1}

	// Setup
	ResetBindings()

	pflags := pflag.NewFlagSet("ApplySpecial", pflag.PanicOnError)
	registerAndBindFlags(flags, pflags)

	args := argsFromFlags("app", flags)
	pflags.Parse(args[1:])

	// Simulate keeping values from the app commandline.
	for _, bf := range GetBindFlags() {
		if bf.Flag.Changed {
			bf.SetValueFrom(bf.Flag)
		}
	}

	Apply()

	pflags = pflag.NewFlagSet("NewFlags", pflag.PanicOnError)
	registerAndBindFlags(flags, pflags)

	// Interactive command line
	argValue := "unique-test-file"
	args = []string{"--" + f1.fk, argValue}
	pflags.Parse(args)

	ApplyFromFlags(pflags)
	v := viper.GetString(f1.bk)
	if v != argValue {
		t.Errorf("failed to set new argvalue from command line. Got: %#v, Expected: %#v",
			v, argValue)
	}

	Apply()

	v = viper.GetString(f1.bk)
	if v != f1.fv {
		t.Errorf("failed to set new argvalue from command line. Got: %#v, Expected: %#v",
			v, f1.fv)
	}

}
func TestMultipleBindsInteractiveLineFlag(t *testing.T) {
	SetDebug(true)
	// Configure
	flags := []flag{f1}

	// Setup
	ResetBindings()

	pflags := pflag.NewFlagSet("ApplySpecial", pflag.PanicOnError)
	registerAndBindFlags(flags, pflags)

	// args := argsFromFlags("app", flags)
	args := []string{"app"}
	pflags.Parse(args[1:])

	// Simulate keeping values from the app commandline.
	for _, bf := range GetBindFlags() {
		if bf.Flag.Changed {
			bf.SetValueFrom(bf.Flag)
		}
	}

	Apply()

	pflags = pflag.NewFlagSet("NewFlags", pflag.PanicOnError)
	registerAndBindFlags(flags, pflags)

	// Interactive command line
	argValue := "unique-test-file"
	args = []string{"--" + f1.fk, argValue}
	pflags.Parse(args)

	ApplyFromFlags(pflags)
	v := viper.GetString(f1.bk)
	if v != argValue {
		t.Errorf("failed to set new argvalue from command line. Got: %#v, Expected: %#v",
			v, argValue)
	}

	Apply()

	// In this case it should be empty as we had no command line value.
	v = viper.GetString(f1.bk)
	if v != "" {
		t.Errorf("failed to set new argvalue from command line. Got: %#v, Expected: %#v",
			v, "")
	}

}
