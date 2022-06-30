package local

import (
	"encoding/json"
	"fmt"

	"github.com/mcdonaldseanp/clibuild/errtype"
	"github.com/mcdonaldseanp/clibuild/validator"
	"github.com/mcdonaldseanp/lookout/localexec"
	"github.com/mcdonaldseanp/lookout/localfile"
	"github.com/mcdonaldseanp/lookout/operation"
	"github.com/mcdonaldseanp/lookout/operparse"
)

func RunAction(actn operation.Action) operation.ActionResult {
	result := operation.ActionResult{
		Action: actn,
	}
	output, logs, cmd_err := localexec.BuildAndRunCommand(actn.Exe, actn.Path, actn.Script, actn.Args)
	if cmd_err != nil {
		result.Succeeded = false
		result.Output = output
		result.Logs = fmt.Sprintf("Error: %s, Logs: %s", cmd_err.(*errtype.ShellError).Message, logs)
	} else {
		result.Succeeded = true
		result.Output = output
		result.Logs = logs
	}
	return result
}

func Run(raw_data []byte, actn_name string) (string, error) {
	err := validator.ValidateParams(fmt.Sprintf(
		`[{"name":"action name","value":"%s","validate":["NotEmpty"]}]`,
		actn_name,
	))
	if err != nil {
		return "", err
	}
	var data operation.Operations
	parse_err := operparse.ParseOperations(raw_data, &data)
	if parse_err != nil {
		return "", parse_err
	}
	actn := operparse.SelectAction(actn_name, data.Actions)
	if actn == nil {
		return "", &errtype.InvalidInput{
			Message: fmt.Sprintf("Name \"%s\" does not match any existing action names", actn_name),
			Origin:  nil,
		}
	}
	result := RunAction(*actn)
	raw_final_result := operation.ActionResults{Actions: make(map[string]operation.ActionResult)}
	raw_final_result.Actions[actn_name] = result
	// The result for actions (for now) is an actionresults set with one action
	// result in the actions field.
	json_output, json_err := json.Marshal(raw_final_result)
	if json_err != nil {
		return "", fmt.Errorf("could not render result as JSON: %s", json_err)
	}
	return string(json_output), nil
}

func CLIRun(maybe_file string, actn_name string) error {
	// ReadFileOrStdin performs validation on maybe_file
	raw_data, err := localfile.ReadFileOrStdin(maybe_file)
	if err != nil {
		return err
	}
	result, err := Run(raw_data, actn_name)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}
