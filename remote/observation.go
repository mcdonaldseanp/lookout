package remote

import (
	"fmt"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/connection"
	"github.com/mcdonaldseanp/lookout/localfile"
)

func Observe(raw_data []byte, username string, target string, port string) (string, error) {
	err := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"username","value":"%s","validate":["NotEmpty"]},
			{"name":"target","value":"%s","validate":["NotEmpty"]},
			{"name":"port","value":"%s","validate":["NotEmpty","IsNumber"]}
		 ]`,
		username,
		target,
		port,
	))
	if err != nil {
		return "", err
	}
	sout, serr, ec, err := connection.RunSSHCommand("$HOME/.lookout/bin/lookout observe local --stdin", string(raw_data), username, target, port)
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

func CLIObserve(maybe_file string, username string, target string, port string) error {
	raw_data, err := localfile.ReadFileOrStdin(maybe_file)
	if err != nil {
		return err
	}
	sout, err := Observe(raw_data, username, target, port)
	if err != nil {
		return err
	}
	fmt.Printf("%s", sout)
	return nil
}
