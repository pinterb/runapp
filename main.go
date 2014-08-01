package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"syscall"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.SetOutput(ioutil.Discard)

	// grab the absolute path of this executable.
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	args := os.Args[1:]

	// basically, we're defining our default behavior here.
	opts := &Options{
		ConfigFilename:         dir + string(os.PathSeparator) + ".env",
		OverrideWithConfigVars: false,
		EnvPrefixFlagName:      "env_prefix",
		OverrideWithEnvVars:    false,
	}
	command, flagettes, err := Parse(opts, args)

	if err != nil {
		log.Fatalf("ParseArgs: %v", err)
	}

	if len(command) == 0 {
		log.Fatalf("No command was parsed(?)")
	}

	if flagettes == nil {
		log.Printf("No command line arguments were parsed(?)")
	}

	// unfortunately, iteration over a map is unordered.  So we're gonna sort a slice...
	finalFlagettes := make(flagetteSlice, 0, len(flagettes))
	for _, value := range flagettes {
		finalFlagettes = append(finalFlagettes, value)
	}
	sort.Sort(finalFlagettes)

	// iterate over our sorted slice of Flagette's to create a new array of (final) arguments
	finalArgs := make([]string, len(finalFlagettes))
	ndx := 0
	for _, value := range finalFlagettes {
		if value.HasEqualAssignmentOper() {
			finalArgs[ndx] = value.Emit()
		} else {
			// for cli arguments passed in without an assignment operator
			finalArgs[ndx] = value.Name()
			copiedArgs := make([]string, len(finalArgs)+1)
			copy(copiedArgs, finalArgs)
			finalArgs = copiedArgs
			ndx++
			finalArgs[ndx] = value.Value()
		}
		ndx++
	}

	cmd := exec.Command(command, finalArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("Running app: ", command, strings.Join(finalArgs, " "))

	if err := cmd.Start(); err != nil {
		log.Fatalf("cmd.Start: %v", err)
	}

	// From http://stackoverflow.com/questions/10385551/get-exit-code-go
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			// The program has exited with an exit code != 0

			// This works on both Unix and Windows. Although package
			// syscall is generally platform dependent, WaitStatus is
			// defined for both Unix and Windows and in both cases has
			// an ExitStatus() method with the same signature.
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				log.Printf("Exit Status: %d", status.ExitStatus())
			}
		} else {
			log.Fatalf("cmd.Wait: %v", err)
		}
	}

	fmt.Println("Finished running app: ", command, strings.Join(finalArgs, " "))

	return 0
}
