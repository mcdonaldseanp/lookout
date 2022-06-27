package remote

import (
	"fmt"

	"github.com/mcdonaldseanp/lookout/connection"
	"github.com/mcdonaldseanp/lookout/localfile"
	"github.com/mcdonaldseanp/lookout/rgerror"
	"github.com/mcdonaldseanp/lookout/validator"
)

func Observe(raw_data []byte, username string, target string, port string) (string, *rgerror.RGerror) {
	rgerr := validator.ValidateParams(fmt.Sprintf(
		`[
			{"name":"username","value":"%s","validate":["NotEmpty"]},
			{"name":"target","value":"%s","validate":["NotEmpty"]},
			{"name":"port","value":"%s","validate":["NotEmpty","IsNumber"]}
		 ]`,
		username,
		target,
		port,
	))
	if rgerr != nil {
		return "", rgerr
	}
	sout, serr, ec, rgerr := connection.RunSSHCommand("$HOME/.lookout/bin/lookout observe local --stdin", string(raw_data), username, target, port)
	if rgerr != nil {
		return sout, &rgerror.RGerror{
			Kind: rgerror.RemoteExecError,
			Message: fmt.Sprintf("lookout client on remote target returned non-zero exit code %d\n\nStdout:\n%s\nStderr:\n%s\n",
				ec,
				sout,
				serr),
			Origin: rgerr.Origin,
		}
	}
	return sout, nil
}

func CLIObserve(maybe_file string, username string, target string, port string) *rgerror.RGerror {
	raw_data, rgerr := localfile.ReadFileOrStdin(maybe_file)
	if rgerr != nil {
		return rgerr
	}
	sout, rgerr := Observe(raw_data, username, target, port)
	if rgerr != nil {
		return rgerr
	}
	fmt.Printf("%s", sout)
	return nil
}
