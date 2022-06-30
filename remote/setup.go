package remote

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/remoteexec"
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
	sout, serr, ec, err := remoteexec.RunSSHCommand(command, "", username, target, port)
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
	json_output, json_err := json.Marshal(output)
	if json_err != nil {
		return fmt.Errorf("could not render result as JSON: %s", json_err)
	}
	fmt.Print(string(json_output))
	return nil
}
