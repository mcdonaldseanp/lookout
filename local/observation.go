package local

import (
	"fmt"
	"strings"

	"github.com/mcdonaldseanp/lookout/localexec"
	"github.com/mcdonaldseanp/lookout/localfile"
	"github.com/mcdonaldseanp/lookout/operation"
	"github.com/mcdonaldseanp/lookout/operparse"
	"github.com/mcdonaldseanp/lookout/render"
)

func RunObservation(name string, obsv operation.Observation, impls map[string]operation.Implement) operation.ObservationResult {
	entity := obsv.Entity
	query := obsv.Query

	for _, impl := range impls {
		if impl.Observes.Query == query && impl.Observes.Entity == entity {
			impl_file := impl.Path
			impl_script := impl.Script
			executable := impl.Exe
			dwld_file, err := DownloadImplement(&impl)
			if err != nil {
				return operation.ObservationResult{
					Succeeded:   false,
					Result:      "failed to download implement",
					Expected:    false,
					Logs:        err.Error(),
					Observation: obsv,
				}
			} else if len(dwld_file) > 0 {
				if len(executable) > 0 {
					impl_file = dwld_file
				} else {
					impl_file = ""
					executable = dwld_file
				}
			}
			args := operparse.ComputeArgs(impl.Observes.Args, obsv)
			output, logs, cmd_err := localexec.BuildAndRunCommand(executable, impl_file, impl_script, args)
			if cmd_err != nil {
				return operation.ObservationResult{
					Succeeded:   false,
					Result:      "error: " + strings.TrimSpace(cmd_err.Error()),
					Expected:    false,
					Logs:        logs,
					Observation: obsv,
				}
			} else {
				result := operation.ObservationResult{
					Succeeded:   true,
					Result:      output,
					Logs:        logs,
					Observation: obsv,
				}
				if obsv.Expect == output || obsv.Expect == "" {
					result.Expected = true
				} else {
					result.Expected = false
				}
				return result
			}
		}
	}
	return operation.ObservationResult{
		Succeeded:   false,
		Result:      "error: No implement found for observation '" + name + "'",
		Observation: obsv,
	}
}

func RunAllObservations(obsvs map[string]operation.Observation, impls map[string]operation.Implement) operation.ObservationResults {
	results := operation.ObservationResults{Observations: make(map[string]operation.ObservationResult)}
	for obsv_name, obsv := range obsvs {
		this_result := RunObservation(obsv_name, obsv, impls)
		results.Observations[obsv_name] = this_result
		results.Total_Observations++
		if this_result.Succeeded == false {
			results.Failed_Observations++
		}
		if this_result.Expected == false {
			results.Unexpected_Observations++
		}
	}
	return results
}

func Observe(raw_data []byte) (string, error) {
	// No validators are required to run here because ParseOperations
	// will use ReadFileOrStdin which performs validation on
	// maybe_file
	var data operation.Operations
	parse_err := operparse.ParseOperations(raw_data, &data)
	if parse_err != nil {
		return "", parse_err
	}
	results := RunAllObservations(data.Observations, data.Implements)
	final_result, parse_err := render.RenderJson(results)
	if parse_err != nil {
		return "", parse_err
	}

	return final_result, nil
}

func CLIObserve(maybe_file string) error {
	// ReadFileOrStdin performs validation on maybe_file
	raw_data, err := localfile.ReadFileOrStdin(maybe_file)
	if err != nil {
		return err
	}
	result, err := Observe(raw_data)
	if err != nil {
		return err
	}
	fmt.Printf(result)
	return nil
}
