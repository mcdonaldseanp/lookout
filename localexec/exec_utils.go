package localexec

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/lookout/localfile"
	"github.com/mcdonaldseanp/lookout/sanitize"
)

func ExecReadOutput(executable string, args []string) (string, string, error) {
	shell_command := exec.Command(executable, args...)
	shell_command.Env = os.Environ()
	var stdout, stderr bytes.Buffer
	shell_command.Stdout = &stdout
	shell_command.Stderr = &stderr
	err := shell_command.Run()
	output := sanitize.ReplaceAllNewlines(stdout.String())
	logs := sanitize.ReplaceAllNewlines(stderr.String())
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
		return "", "", fmt.Errorf("Could not create tmp file!")
	}
	filename := f.Name()
	defer os.Remove(filename) // clean up
	localfile.OverwriteFile(filename, []byte(script))
	final_args := append([]string{filename}, args...)
	return ExecReadOutput(executable, final_args)
}

func BuildAndRunCommand(executable string, file string, script string, args []string) (string, string, error) {
	var output, logs string
	var err error
	if len(file) > 0 {
		final_args := append([]string{file}, args...)
		output, logs, err = ExecReadOutput(executable, final_args)
	} else if len(script) > 0 {
		output, logs, err = ExecScriptReadOutput(executable, script, args)
	} else {
		output, logs, err = ExecReadOutput(executable, args)
	}
	if err != nil {
		return output, logs, err
	}

	return output, logs, nil
}
