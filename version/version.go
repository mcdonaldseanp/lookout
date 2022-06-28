package version

import "fmt"

const VERSION string = "v0.0.0"

func ReleaseArtifact(name string) string {
	return fmt.Sprintf("https://github.com/mcdonaldseanp/lookout/releases/download/%s/%s", VERSION, name)
}
