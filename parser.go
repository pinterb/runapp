package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	ini "github.com/rakyll/goini"
)

// A Flagette represents the state of a flag. It was pulled from https://github.com/spf13/pflag/blob/master/flag.go
type Flagette struct {
	name       string // name as it appears on command line
	assignOper string // assignment operator
	value      string // value as set
}

// Value is the interface to the dynamic value stored in a flag.
// (The default value is represented as a string.)
type Value interface {
	String() string
	Set(string) error
	Type() string
}

// Options for defining location of configuration file and/or environment variables prefix
type Options struct {
	ConfigFilename         string // absolute path to application's configuration file
	OverrideWithConfigVars bool   // override any existing Flagette's matching config file variable
	EnvPrefixFlagName      string // the flag used to specify the application prefix used to identify environment variables
	OverrideWithEnvVars    bool   // override any existing Flagette's matching environment variable
}

// SetName will set the name of the flagette.
func (f *Flagette) SetName(name string) {
	f.name = name
}

// Name will return the name of the flagette.
func (f Flagette) Name() string {
	return f.name
}

// SetAssignOper will set the assignment operator of the flagette.
func (f *Flagette) SetAssignOper(value string) {
	fmt.Printf("Setting Flag AssignOper: %s\n", value)
	f.assignOper = value
}

// AssignOper will return the assignment operator of the flagette.
func (f Flagette) AssignOper() string {
	return f.assignOper
}

// HasEqualAssignmentOper will indicate if the assignment operator used was an equals sign (eg "="); If false, then a " " was used.
func (f Flagette) HasEqualAssignmentOper() bool {
	return f.assignOper == "="
}

// SetValue will set the value of the flagette.
func (f *Flagette) SetValue(value string) {
	fmt.Printf("Setting Flag Value: %s\n", value)
	f.value = value
}

// Value will return the value of the flagette.
func (f Flagette) Value() string {
	return f.value
}

// HasValue will indicate if the value has been set.
func (f Flagette) HasValue() bool {
	return len(f.value) > 0
}

// Emit will return the current name/assignOper/value combination.
func (f Flagette) Emit() string {
	return f.name + f.assignOper + f.value
}

// Parse will take some command line arguments and parse them out into a struct that we can work with.
func Parse(opts *Options, args []string) (c string, f map[string]*Flagette, e error) {
	fmt.Println(fmt.Sprintf("ParseArgs: %s", "Starting..."))

	// Define our return values
	var command string
	var newArgs map[string]*Flagette
	var err error

	if len(args) > 0 {
		first := args[0]
		args = args[1:]
		//	newArgs = make(map[string]*Flagette)

		// if first arg is a flag (either long or short) and not --help | -h | -v | --version then it's probably an error
		if first == "-h" || first == "--help" {
			command = "help"
			return command, newArgs, err

		} else if first == "-v" || first == "--version" || first == "-h" || first == "--help" {
			command = "version"
			return command, newArgs, err

		} else if strings.HasPrefix(first, "--") || strings.HasPrefix(first, "-") {
			return command, newArgs, errors.New("The command, '" + first + "' appears to be a flag of some sort.")

		}

		// some variable initialization
		command = first
		newArgs, err = parseCommandLineArgs(args)
	}

	fmt.Println("")
	fmt.Println("Flagettes: after command line parse")
	for key, value := range newArgs {
		fmt.Println("Key:", key, " --> Name:", value.Name(), "AssignOper:", value.AssignOper(), "Value:", value.Value(), " (", value.Emit(), ")")
	}

	newArgs, err = parseEnvVars(opts, newArgs)

	fmt.Println("")
	fmt.Println("Flagettes: after env vars parse")
	for key, value := range newArgs {
		fmt.Println("Key:", key, " --> Name:", value.Name(), "AssignOper:", value.AssignOper(), "Value:", value.Value(), " (", value.Emit(), ")")
	}

	fmt.Println(fmt.Sprintf("ParseArgs: %s", "Ending..."))
	return command, newArgs, err
}

// Parse command line arguments passed in by the user.
func parseCommandLineArgs(args []string) (f map[string]*Flagette, e error) {
	fmt.Println(fmt.Sprintf("parseCommandLineArgs: %s", "Starting..."))

	// Define our return values
	var newArgs map[string]*Flagette
	var err error

	if len(args) > 0 {
		//	first := args[0]
		//	args = args[1:]
		newArgs = make(map[string]*Flagette)

		// some variable initialization
		currentFlg := new(Flagette)
		currentFlg.SetName("")

		for _, arg := range args {

			fmt.Println("")
			fmt.Println("")
			fmt.Printf("Current Flag: %s\n", currentFlg.Name())
			fmt.Printf("Current Flag Value: %s\n", currentFlg.Value())

			if strings.HasPrefix(arg, "--") {
				fmt.Printf("We have a long flag: %s\n", arg)
				currentFlg = new(Flagette)
				strippedLong := strings.TrimPrefix(arg, "--")
				var lMapKey string
				if strings.Contains(strippedLong, "=") {
					longSplits := strings.Split(strippedLong, "=")
					lMapKey = longSplits[0]
					currentFlg.SetName("--" + longSplits[0])
					currentFlg.SetAssignOper("=")
					currentFlg.SetValue(strings.Join(longSplits[1:], ""))
				} else {
					lMapKey = strippedLong
					currentFlg.SetName("--" + strippedLong)
					currentFlg.SetAssignOper(" ")
				}

				newArgs[lMapKey] = currentFlg

			} else if strings.HasPrefix(arg, "-") {
				fmt.Printf("We have a short flag: %s\n", arg)
				currentFlg = new(Flagette)
				strippedShort := strings.TrimPrefix(arg, "-")
				var sMapKey string
				if strings.Contains(strippedShort, "=") {
					shortSplits := strings.Split(strippedShort, "=")
					sMapKey = shortSplits[0]
					currentFlg.SetName("-" + shortSplits[0])
					currentFlg.SetAssignOper("=")
					currentFlg.SetValue(strings.Join(shortSplits[1:], ""))
				} else {
					sMapKey = strippedShort
					currentFlg.SetName("-" + strippedShort)
					currentFlg.SetAssignOper(" ")
				}

				newArgs[sMapKey] = currentFlg

			} else if currentFlg.Name() != "" && !currentFlg.HasValue() {
				fmt.Printf("We have a a value: %s\n", arg)
				currentFlg.SetValue(arg)

			} else {
				fmt.Printf("I don't know what we have: %s\n", arg)
			}
		}
	}

	fmt.Println(fmt.Sprintf("parseCommandLineArgs: %s", "Ending..."))
	return newArgs, err
}

// Parse any environment variables found matching application prefix.
func parseEnvVars(opts *Options, flagettes map[string]*Flagette) (f map[string]*Flagette, err error) {
	fmt.Println(fmt.Sprintf("parseEnvVars: %s", "Starting..."))
	if opts.EnvPrefixFlagName != "" {
		if appKeyFlagette, ok := flagettes[opts.EnvPrefixFlagName]; ok {
			envVars, err := loadEnvVars()
			if err == nil {
				appKey := strings.ToUpper(appKeyFlagette.Value())
				delete(flagettes, opts.EnvPrefixFlagName)
				if !strings.HasSuffix(appKey, "_") {
					appKey += "_"
				}
				for envKey, envValue := range envVars {
					if strings.HasPrefix(envKey, appKey) {
						lookupKey := strings.ToLower(strings.TrimPrefix(envKey, appKey))
						fmt.Println("LookupKEY ", lookupKey)
						if _, ok := flagettes[lookupKey]; ok {
							if !opts.OverrideWithEnvVars {
								fmt.Println(lookupKey, ": Was already entered at the command line")
								continue
							}
						}
						fmt.Println("Putting ", lookupKey)
						flagettes[lookupKey] = &Flagette{
							name:       "--" + lookupKey,
							assignOper: "=",
							value:      envValue,
						}
					}
				}
			}
		}
	}

	fmt.Println(fmt.Sprintf("parseEnvVars: %s", "Ending..."))

	fmt.Println("")
	fmt.Println("Flagettes: after env var parse")
	for key, value := range flagettes {
		fmt.Println("Key:", key, " --> Name:", value.Name(), "AssignOper:", value.AssignOper(), "Value:", value.Value(), " (", value.Emit(), ")")
	}

	return flagettes, nil
}

// Load all environment variables.  See https://coderwall.com/p/kjuyqw
func loadEnvVars() (env map[string]string, err error) {
	fmt.Println(fmt.Sprintf("loadEnvVars: %s", "Starting..."))

	getenvironment := func(data []string, getkeyval func(item string) (key, val string)) map[string]string {
		items := make(map[string]string)
		for _, item := range data {
			key, val := getkeyval(item)
			items[key] = val
		}
		return items
	}

	environment := getenvironment(os.Environ(), func(item string) (key, val string) {
		splits := strings.Split(item, "=")
		key = splits[0]
		val = strings.Join(splits[1:], "=")
		return
	})

	//fmt.Println(environment["PATH"])

	fmt.Println(fmt.Sprintf("loadEnvVars: %s", "Ending..."))

	/* for key, value := range environment {
		fmt.Println("Key:", key, "Value:", value)
	} */
	return environment, nil
}

//  Load any config file variables.
func loadConfigFileVars(opts *Options, flagettes map[string]*Flagette) (err error) {
	fmt.Println(fmt.Sprintf("loadConfigFileVars: %s", "Starting..."))
	fmt.Println(fmt.Sprintf("loadConfigFileVars: %s", "Ending..."))
	return nil
}

// Load configuration file
func readConfig(opts *Options) (d *ini.Dict, err error) {
	fmt.Println(fmt.Sprintf("loadConfig: %s", "Starting..."))
	var dict ini.Dict
	if opts.ConfigFilename != "" {
		dict, err = ini.Load(opts.ConfigFilename)
		if err != nil {
			return nil, err
		}
	} else {
		dict = make(ini.Dict, 0)
	}
	fmt.Println(fmt.Sprintf("loadConfig: %s", "Ending..."))
	return &dict, nil
}
