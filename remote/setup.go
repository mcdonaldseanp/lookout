package remote

import (
	"fmt"
	"strings"

	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/connection"
	"github.com/mcdonaldseanp/lookout/render"
	"github.com/mcdonaldseanp/lookout/rgerror"
	"github.com/mcdonaldseanp/lookout/version"
)

func Setup(username string, target string, port string) (string, string, error) {
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
		return "", "", rgerr
	}
	command := fmt.Sprintf(
		`#!/usr/bin/env bash

		mkdir -p $HOME/.lookout/bin 1>&2
		curl -L %s > $HOME/.lookout/bin/lookout
		chmod 755 $HOME/.lookout/bin/lookout 1>&2`,
		version.ReleaseArtifact("lookout"),
	)
	sout, serr, ec, rgerr := connection.RunSSHCommand(command, "", username, target, port)
	if rgerr != nil {
		return "", "", &rgerror.RGerror{
			Kind: rgerror.RemoteExecError,
			Message: fmt.Sprintf("lookout client on remote target returned non-zero exit code %d\n\nStdout:\n%s\nStderr:\n%s\n",
				ec,
				sout,
				serr),
			Origin: rgerr.(*rgerror.RGerror).Origin,
		}
	}
	return sout, serr, nil
}

func CLISetup(username string, target string, port string) error {
	_, serr, rgerr := Setup(username, target, port)
	if rgerr != nil {
		return rgerr
	}
	output := make(map[string]interface{})
	output["ok"] = true
	output["logs"] = strings.TrimSpace(serr)
	final_result, json_rgerr := render.RenderJson(output)
	if json_rgerr != nil {
		return json_rgerr
	}
	fmt.Printf(final_result)
	return nil
}
