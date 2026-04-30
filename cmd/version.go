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
	var v struct {
		Version string `json:"version"`
	}
	_ = json.Unmarshal(versionJson, &v)

	var commit string
	var timestamp time.Time
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, kv := range info.Settings {
			switch kv.Key {
			case "vcs.revision":
				commit = kv.Value
			case "vcs.time":
				timestamp, _ = time.Parse("2006-01-02T15:04:05Z", kv.Value)
			}
		}
	}

	shortCommit := commit
	if len(shortCommit) > 10 {
		shortCommit = shortCommit[:10]
	}

	fmt.Printf("%s %s %s\n", v.Version, shortCommit, timestamp.Format(time.DateTime))
	return nil
}
