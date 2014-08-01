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
	sequence   int    // this struct needs to be sortable
	name       string // name as it appears on command line
	assignOper string // assignment operator
	value      string // value as set
}

// Options for defining location of configuration file and/or environment variables prefix
type Options struct {
	ConfigFilename         string // absolute path to application's configuration file
	OverrideWithConfigVars bool   // override any existing Flagette's matching config file variable
	EnvPrefixFlagName      string // the flag used to specify the application prefix used to identify environment variables
	OverrideWithEnvVars    bool   // override any existing Flagette's matching environment variable
}

// SetSequence will set the sequence of the flagette.
func (f *Flagette) SetSequence(seq int) {
	f.sequence = seq
}

// Sequence will return the sequence of the flagette.
func (f Flagette) Sequence() int {
	return f.sequence
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

// This type is used for sorting Flagette's
type flagetteSlice []*Flagette

// Len is part of sort.Interface
func (f flagetteSlice) Len() int {
	return len(f)
}

// Swap is part of sort.Interface
func (f flagetteSlice) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

// Less is part of sort.Interface. We sort on Flagette's sequence value.
func (f flagetteSlice) Less(i, j int) bool {
	return f[i].Sequence() < f[j].Sequence()
}

// Parse will take some command line arguments and parse them out into a struct that we can work with.
func Parse(opts *Options, args []string) (c string, f map[string]*Flagette, e error) {

	// define our return values
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

	err = parseEnvVars(opts, newArgs)
	if err != nil {
		fmt.Println("Flagettes: back from env vars parse. But recieved error: ", err)
	}

	err = parseConfigFileVars(opts, newArgs)
	if err != nil {
		fmt.Println("Flagettes: back from config file parse. But recieved error: ", err)
	}

	return command, newArgs, err
}

// Parse command line arguments passed in by the user.
func parseCommandLineArgs(args []string) (f map[string]*Flagette, e error) {

	// define our return values
	var newArgs map[string]*Flagette
	var err error

	// define our counter
	counter := 0

	if len(args) > 0 {
		newArgs = make(map[string]*Flagette)

		// some variable initialization
		currentFlg := new(Flagette)
		currentFlg.SetName("")

		for _, arg := range args {

			if strings.HasPrefix(arg, "--") {
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

				counter++
				currentFlg.SetSequence(counter)
				newArgs[lMapKey] = currentFlg

			} else if strings.HasPrefix(arg, "-") {
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

				counter++
				currentFlg.SetSequence(counter)
				newArgs[sMapKey] = currentFlg

			} else if currentFlg.Name() != "" && !currentFlg.HasValue() {
				currentFlg.SetValue(arg)

			} else {
				fmt.Printf("I don't know what we have: %s\n", arg)
			}
		}
	}

	return newArgs, err
}

// Parse any environment variables found matching application prefix.
func parseEnvVars(opts *Options, flagettes map[string]*Flagette) (e error) {
	if opts.EnvPrefixFlagName != "" {
		if appKeyFlagette, ok := flagettes[opts.EnvPrefixFlagName]; ok {
			envVars, err := loadEnvVars()
			if err != nil {
				return err
			}
			counter := len(flagettes)
			appKey := strings.ToUpper(appKeyFlagette.Value())
			delete(flagettes, opts.EnvPrefixFlagName)
			if !strings.HasSuffix(appKey, "_") {
				appKey += "_"
			}
			for envKey, envValue := range envVars {
				if strings.HasPrefix(envKey, appKey) {
					lookupKey := strings.ToLower(strings.TrimPrefix(envKey, appKey))
					if _, ok := flagettes[lookupKey]; ok {
						if !opts.OverrideWithEnvVars {
							continue
						}
					}

					counter++
					flagettes[lookupKey] = &Flagette{
						name:       "--" + lookupKey,
						assignOper: "=",
						value:      envValue,
						sequence:   counter,
					}
				}
			}
		}
	}

	return nil
}

// Load all environment variables.  See https://coderwall.com/p/kjuyqw
func loadEnvVars() (env map[string]string, err error) {

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

	return environment, nil
}

// Parse any config file variables.
func parseConfigFileVars(opts *Options, flagettes map[string]*Flagette) (e error) {
	config, err := readConfig(opts)
	if err != nil {
		if strings.Contains(err.Error(), "no such file or directory") {
			return nil
		}
		return err
	}

	counter := len(flagettes)
	// ini files can have sections.
	for configSection, configDetails := range config {
		// inside a section
		for configName, configValue := range configDetails {

			lookupKey := strings.ToLower(configName)
			if len(configSection) > 0 {
				lookupKey = strings.ToLower(configSection + "_" + configName)
			}

			if _, ok := flagettes[lookupKey]; ok {
				if !opts.OverrideWithConfigVars {
					continue
				}
			}

			counter++
			flagettes[lookupKey] = &Flagette{
				name:       "--" + lookupKey,
				assignOper: "=",
				value:      configValue,
				sequence:   counter,
			}

		}
	}

	return err
}

// Load configuration file
func readConfig(opts *Options) (d ini.Dict, err error) {
	var dict ini.Dict
	if opts.ConfigFilename != "" {
		dict, err = ini.Load(opts.ConfigFilename)
		if err != nil {
			return nil, err
		}
	} else {
		dict = make(ini.Dict, 0)
	}
	return dict, nil
}
