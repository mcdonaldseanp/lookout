package render

import (
	"encoding/json"
	"fmt"

	"github.com/mcdonaldseanp/lookout/rgerror"
)

func RenderJson(data interface{}) (string, *rgerror.RGerror) {
	json_output, json_err := json.Marshal(data)
	if json_err != nil {
		return "", &rgerror.RGerror{
			Kind:    rgerror.ExecError,
			Message: fmt.Sprintf("Could not render result as JSON: %s\n", json_err),
			Origin:  json_err,
		}
	}
	return string(json_output), nil
}
