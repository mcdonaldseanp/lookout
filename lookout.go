package main

import (
	"flag"
	"os"

	"github.com/mcdonaldseanp/clibuild/cli"
	"github.com/mcdonaldseanp/lookout/local"
	"github.com/mcdonaldseanp/lookout/localdata"
	"github.com/mcdonaldseanp/lookout/remote"
	"github.com/mcdonaldseanp/lookout/version"
)

func main() {
	// Use flagsets from the https://pkg.go.dev/flag package
	// to define CLI flags
	//
	// None of the commands below should call .Parse on any
	// flagsets directly. cli.ShouldHaveArgs() will call .Parse
	// on the flagset if it is passed one.
	//
	// Things need to be parsed inside cli.ShouldHaveArgs so that
	// the flag package can ignore any required commands
	// before parsing
	local_flag_set := flag.NewFlagSet("local_options", flag.ExitOnError)
	local_input_file := local_flag_set.String("file", "", "Path to spec yaml file (must use one of --file or --stdin)")
	local_use_stdin := local_flag_set.Bool("stdin", false, "Read spec from stdin (must use one of --file or --stdin)")

	remote_flag_set := flag.NewFlagSet("remote_options", flag.ExitOnError)
	remote_input_file := remote_flag_set.String("file", "", "Path to spec yaml file (must use one of --file or --stdin)")
	remote_use_stdin := remote_flag_set.Bool("stdin", false, "Read spec from stdin (must use one of --file or --stdin)")
	username := remote_flag_set.String("user", os.Getenv("USER"), "Username to use when connecting via SSH")
	port := remote_flag_set.String("port", "22", "Port to use for ssh connections")

	setup_flag_set := flag.NewFlagSet("setup_options", flag.ExitOnError)
	setup_username := setup_flag_set.String("user", os.Getenv("USER"), "Username to use when connecting via SSH")
	setup_port := setup_flag_set.String("port", "22", "Port to use for ssh connections")

	// All CLI commands should follow naming rules of powershell approved verbs:
	// https://docs.microsoft.com/en-us/powershell/scripting/developer/cmdlet/approved-verbs-for-windows-powershell-commands?view=powershell-7.2
	//
	// Also, try to keep these in alphabetical order. The list is already long enough
	command_list := []cli.Command{
		{
			Verb:     "observe",
			Noun:     "local",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "lookout observe local [FLAGS]"
				description := "Run observation code on the local system and print out the resulting observations"
				cli.ShouldHaveArgs(0, usage, description, local_flag_set)
				input_file, err := localdata.ChooseFileOrStdin(*local_input_file, *local_use_stdin)
				if err != nil {
					cli.HandleCommandError(err, usage, description, local_flag_set)
				}
				cli.HandleCommandError(
					local.CLIObserve(input_file),
					usage,
					description,
					local_flag_set,
				)
			},
		},
		{
			Verb:     "observe",
			Noun:     "remote",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "lookout observe remote [TARGET] [FLAGS]"
				description := "Run observation on a target"
				cli.ShouldHaveArgs(1, usage, description, remote_flag_set)
				input_file, err := localdata.ChooseFileOrStdin(*remote_input_file, *remote_use_stdin)
				if err != nil {
					cli.HandleCommandError(err, usage, description, remote_flag_set)
				}
				cli.HandleCommandError(
					remote.CLIObserve(input_file, *username, os.Args[3], *port),
					usage,
					description,
					remote_flag_set,
				)
			},
		},
		{
			Verb:     "react",
			Noun:     "local",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "lookout react local [FLAGS]"
				description := "React to an observation on the local system"
				cli.ShouldHaveArgs(0, usage, description, local_flag_set)
				input_file, err := localdata.ChooseFileOrStdin(*local_input_file, *local_use_stdin)
				if err != nil {
					cli.HandleCommandError(err, usage, description, local_flag_set)
				}
				cli.HandleCommandError(
					local.CLIReact(input_file),
					usage,
					description,
					local_flag_set,
				)
			},
		},
		{
			Verb:     "react",
			Noun:     "remote",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "lookout react remote [TARGET] [FLAGS]"
				description := "React to an observation on a target"
				cli.ShouldHaveArgs(1, usage, description, remote_flag_set)
				input_file, err := localdata.ChooseFileOrStdin(*remote_input_file, *remote_use_stdin)
				if err != nil {
					cli.HandleCommandError(err, usage, description, remote_flag_set)
				}
				cli.HandleCommandError(
					remote.CLIReact(input_file, *username, os.Args[3], *port),
					usage,
					description,
					remote_flag_set,
				)
			},
		},
		{
			Verb:     "run",
			Noun:     "local",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "lookout run local [ACTION NAME] [FLAGS]"
				description := "Run an action on the local system"
				cli.ShouldHaveArgs(1, usage, description, local_flag_set)
				input_file, err := localdata.ChooseFileOrStdin(*local_input_file, *local_use_stdin)
				if err != nil {
					cli.HandleCommandError(err, usage, description, local_flag_set)
				}
				cli.HandleCommandError(
					local.CLIRun(input_file, os.Args[3]),
					usage,
					description,
					local_flag_set,
				)
			},
		},
		{
			Verb:     "run",
			Noun:     "remote",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "lookout run remote [ACTION NAME] [TARGET] [FLAGS]"
				description := "Run actions on a target"
				cli.ShouldHaveArgs(2, usage, description, remote_flag_set)
				input_file, err := localdata.ChooseFileOrStdin(*remote_input_file, *remote_use_stdin)
				if err != nil {
					cli.HandleCommandError(err, usage, description, remote_flag_set)
				}
				cli.HandleCommandError(
					remote.CLIRun(input_file, os.Args[3], *username, os.Args[4], *port),
					usage,
					description,
					remote_flag_set,
				)
			},
		},
		{
			Verb:     "setup",
			Noun:     "remote",
			Supports: []string{"linux", "windows"},
			ExecutionFn: func() {
				usage := "lookout setup remote [TARGET] [FLAGS]"
				description := "Run actions on a target"
				cli.ShouldHaveArgs(1, usage, description, setup_flag_set)
				cli.HandleCommandError(
					remote.CLISetup(*setup_username, os.Args[3], *setup_port),
					usage,
					description,
					setup_flag_set,
				)
			},
		},
	}

	cli.RunCommand("lookout", version.VERSION, command_list)
}
