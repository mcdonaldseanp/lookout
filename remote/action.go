package remote

import (
	"fmt"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/localdata"
	"github.com/mcdonaldseanp/lookout/remoteexec"
)

func Run(raw_data []byte, actn_name string, username string, target string, port string) (string, error) {
	err := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"action name","value":"%s","validate":["NotEmpty"]},
			{"name":"username","value":"%s","validate":["NotEmpty"]},
			{"name":"target","value":"%s","validate":["NotEmpty"]},
			{"name":"port","value":"%s","validate":["NotEmpty","IsNumber"]}
		 ]`,
		actn_name,
		username,
		target,
		port,
	))
	if err != nil {
		return "", err
	}
	command := fmt.Sprintf("$HOME/.lookout/bin/lookout run local \"%s\" --stdin", actn_name)
	sout, serr, ec, err := remoteexec.RunSSHCommand(command, string(raw_data), username, target, port)
	if err != nil {
		origin := err
		if errtype_origin, ok := origin.(*errtype.RemoteShellError); ok {
			origin = errtype_origin.Origin
		}
		return sout, &errtype.RemoteShellError{
			Message: fmt.Sprintf("lookout client on remote target returned non-zero exit code %d\n\nStdout:\n%s\nStderr:\n%s\n",
				ec,
				sout,
				serr),
			Origin: origin,
		}
	}
	return sout, nil
}

func CLIRun(maybe_file string, actn_name string, username string, target string, port string) error {
	raw_data, err := localdata.ReadFileOrStdin(maybe_file)
	if err != nil {
		return err
	}
	sout, err := Run(raw_data, actn_name, username, target, port)
	if err != nil {
		return err
	}
	fmt.Printf("%s", sout)
	return nil
}
