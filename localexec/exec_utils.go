package localexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/lookout/localdata"
	"github.com/mcdonaldseanp/lookout/sanitize"
)

func ExecReadOutput(command_string string, args ...string) (string, string, error) {
	if runtime.GOOS == "linux" && isWinPath(command_string) {
		translated_cmd, err := wslPathConvert(command_string)
		if err != nil {
			return "", "", err
		}
		command_string = translated_cmd
	}
	shell_command := exec.Command(command_string, args...)
	shell_command.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	shell_command.Stdout = &stdout
	shell_command.Stderr = &stderr
	err := shell_command.Run()
	output := stdout.String()
	logs := stderr.String()
	if err != nil {
		return output, logs, &errtype.ShellError{
			Message: fmt.Sprintf("Command '%s' failed:\n%s\nstderr:\n%s", shell_command, err, logs),
			Origin:  err,
		}
	}
	return output, logs, nil
}

func ExecScriptReadOutput(executable string, script string, args []string) (string, string, error) {
	f, err := os.CreateTemp("", "lookout_script")
	if err != nil {
		return "", "", fmt.Errorf("could not create tmp file")
	}
	filename := f.Name()
	defer os.Remove(filename) // clean up
	localdata.OverwriteFile(filename, []byte(script))
	final_args := append([]string{filename}, args...)
	return ExecReadOutput(executable, final_args...)
}

func BuildAndRunCommand(executable string, file string, script string, args []string) (string, string, error) {
	var output, logs string
	var err error
	if len(file) > 0 {
		final_args := append([]string{file}, args...)
		output, logs, err = ExecReadOutput(executable, final_args...)
	} else if len(script) > 0 {
		output, logs, err = ExecScriptReadOutput(executable, script, args)
	} else {
		output, logs, err = ExecReadOutput(executable, args...)
	}
	if err != nil {
		return output, logs, err
	}

	return output, logs, nil
}

// ExecAsShell always writes everything to stderr so that
// any resulting functionality can return something useful
// to the CLI
func ExecAsShell(command_string string, args ...string) error {
	if runtime.GOOS == "linux" && isWinPath(command_string) {
		translated_cmd, err := wslPathConvert(command_string)
		if err != nil {
			return err
		}
		command_string = translated_cmd
	}
	shell_command := exec.Command(command_string, args...)
	shell_command.Env = os.Environ()
	shell_command.Stdout = os.Stderr
	shell_command.Stderr = os.Stderr
	shell_command.Stdin = os.Stdin
	err := shell_command.Run()
	if err != nil {
		return &errtype.ShellError{
			Message: fmt.Sprintf("command '%s' failed:\n%s", shell_command, err),
			Origin:  err,
		}
	}
	return nil
}

func ExecDetached(command_string string, args ...string) (*exec.Cmd, error) {
	if runtime.GOOS == "linux" && isWinPath(command_string) {
		translated_cmd, err := wslPathConvert(command_string)
		if err != nil {
			return nil, err
		}
		command_string = translated_cmd
	}
	shell_command := exec.Command(command_string, args...)
	shell_command.Env = os.Environ()
	err := shell_command.Start()
	if err != nil {
		return nil, &errtype.ShellError{
			Message: fmt.Sprintf("Command '%s' failed to start:\n%s", shell_command, err),
			Origin:  err,
		}
	}
	return shell_command, nil
}

func wslPathConvert(command_string string) (string, error) {
	if runtime.GOOS == "linux" && isWinPath(command_string) {
		wsl_path, err_log, err := ExecReadOutput("wslpath", "-u", command_string)
		if err != nil {
			return "", fmt.Errorf("detected use of windows path on linux, failed attempt to use wslpath: %s\ntrace: %s", err, err_log)
		}
		command_string = sanitize.ReplaceAllNewlines(wsl_path)
	}
	return command_string, nil
}

func isWinPath(command_string string) bool {
	return strings.HasPrefix(command_string, "C:\\") || strings.HasPrefix(command_string, "C:/")
}
