package cmd

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"runtime/debug"
	"time"

	"github.com/urfave/cli/v2"
)

//go:embed version.json
var versionJson []byte

func Version(cCtx *cli.Context) error {
	var commit string
	var timestamp time.Time

	info, _ := debug.ReadBuildInfo()
	for _, kv := range info.Settings {
		switch kv.Key {
		case "vcs.revision":
			commit = kv.Value
		case "vcs.time":
			timestamp, _ = time.Parse("2006-01-02T15:04:05Z", kv.Value)
		}
	}
	var versionMap map[string]string
	_ = json.Unmarshal(versionJson, &versionMap)

	fmt.Printf("%s %s %s\n", versionMap["message"], commit[:10], timestamp.Format(time.DateTime))
	return nil
}
