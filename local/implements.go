package local

import (
	"os"

	"github.com/mcdonaldseanp/lookout/localdata"
	"github.com/mcdonaldseanp/lookout/operation"
	"github.com/mcdonaldseanp/lookout/remotedata"
)

const IMPLS_LOC string = ".lookout/impls"

func DownloadImplement(impl *operation.Implement) (string, error) {
	if len(impl.Source_Url) < 1 || len(impl.Source_File) < 1 {
		return "", nil
	}
	raw_data, err := remotedata.Download(impl.Source_Url)
	if err != nil {
		return "", err
	}

	file_loc := os.Getenv("HOME") + "/" + IMPLS_LOC + "/" + impl.Source_File
	return file_loc, localdata.OverwriteFile(file_loc, raw_data)
}
