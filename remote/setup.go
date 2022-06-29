package remote

import (
	"fmt"
	"strings"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/connection"
	"github.com/mcdonaldseanp/lookout/render"
	"github.com/mcdonaldseanp/lookout/version"
)

func Setup(username string, target string, port string) (string, string, error) {
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
		return "", "", err
	}
	command := fmt.Sprintf(
		`#!/usr/bin/env bash

		mkdir -p $HOME/.lookout/bin 1>&2
		mkdir -p $HOME/.lookout/impls 1>&2
		curl -L %s > $HOME/.lookout/bin/lookout
		chmod 755 $HOME/.lookout/bin/lookout 1>&2`,
		version.ReleaseArtifact("lookout"),
	)
	sout, serr, ec, err := connection.RunSSHCommand(command, "", username, target, port)
	if err != nil {
		origin := err
		if errtype_origin, ok := origin.(*errtype.RemoteShellError); ok {
			origin = errtype_origin.Origin
		}
		return sout, serr, &errtype.RemoteShellError{
			Message: fmt.Sprintf("attempt to download lookout client on remote target returned non-zero exit code %d\n\nStdout:\n%s\nStderr:\n%s\n",
				ec,
				sout,
				serr),
			Origin: origin,
		}
	}
	return sout, serr, nil
}

func CLISetup(username string, target string, port string) error {
	_, serr, err := Setup(username, target, port)
	if err != nil {
		return err
	}
	output := make(map[string]interface{})
	output["ok"] = true
	output["logs"] = strings.TrimSpace(serr)
	final_result, json_err := render.RenderJson(output)
	if json_err != nil {
		return json_err
	}
	fmt.Printf(final_result)
	return nil
}
