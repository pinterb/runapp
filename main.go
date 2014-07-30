package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

func main() {
	os.Exit(realMain())
}

func realMain() int {
	log.SetOutput(ioutil.Discard)

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	args := os.Args[1:]

	// Get the command line args. We shortcut "--version" and "-v" to
	// just show the version.
	/*	args := os.Args[1:]
		for _, arg := range args {
			if arg == "-v" || arg == "--version" {
				newArgs := make([]string, len(args)+1)
				newArgs[0] = "version"
				copy(newArgs[1:], args)
				args = newArgs
				break
			}
			fmt.Println("Arg: " + arg)
		} */

	origArgsInput := strings.Join(args[:], " ")
	//finalArgs := args[1:]

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
		log.Printf("No command was parsed(?)")
	} else {
		fmt.Println("")
		fmt.Println("Command: ", command)
	}

	if flagettes == nil {
		log.Printf("No command line arguments were parsed(?)")
		/*	} else {
			fmt.Println("")
			fmt.Println("FINAL FLAGS")
			fmt.Println("")
			for key, value := range flagettes {
				fmt.Println("Key:", key, " --> Name:", value.Name(), "AssignOper:", value.AssignOper(), "Value:", value.Value(), " (", value.Emit(), ")")
			} */
	}

	//finalArgsAsString := ""
	finalArgs := make([]string, len(flagettes))
	ndx := 0
	for _, value := range flagettes {
		if value.HasEqualAssignmentOper() {
			//finalArgs[ndx] = strings.Trim(value.Emit(), " ")
			finalArgs[ndx] = value.Emit()
		} else {
			finalArgs[ndx] = value.Name()
			copiedArgs := make([]string, len(finalArgs)+1)
			copy(copiedArgs, finalArgs)
			finalArgs = copiedArgs
			ndx++
			finalArgs[ndx] = value.Value()
		}
		ndx++
	}
	/*
		for meh := range finalArgs {
			fmt.Println("==>", finalArgs[meh], "<==")
		}*/

	//fmt.Println("feh: " + finalArgsAsString)

	//cmd := exec.Command(args[0], finalArgs...)
	//cmd := exec.Command(command, finalArgsAsString)
	cmd := exec.Command(command, finalArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	fmt.Println("feh: " + strings.Join(finalArgs, " "))
	fmt.Println("Executing command: " + origArgsInput)

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

	//err := cmd.Run()
	fmt.Println("Finished running command: " + origArgsInput)
	/*	command, flagettes, err := ParseArgs(args)

		if err != nil {
			log.Fatalf("ParseArgs: %v", err)
		}

		if len(command) == 0 {
			log.Printf("No command was parsed(?)")
		} else {
			fmt.Println("")
			fmt.Println("Command: ", command)
		}

		if flagettes == nil {
			log.Printf("No command line arguments were parsed(?)")
		} else {
			fmt.Println("")
			fmt.Println("")
			for key, value := range flagettes {
				fmt.Println("Key:", key, " --> Name:", value.Name(), "AssignOper:", value.AssignOper(), "Value:", value.Value(), " (", value.Emit(), ")")
			}
		} */

	return 0
}
