package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"runtime/debug"
)

// Default values.
// These variables can be overridden by -ldflags.
var (
	Version = "dev"
	Commit  = "none"
)

type BuildInfo struct {
	Version   string `json:"version"`
	Commit    string `json:"commit"`
	GoVersion string `json:"go_version"`
	Modified  string `json:"modified"`
}

// Get build information from the running binary.
func versionInfo() BuildInfo {
	result := BuildInfo{
		Version: Version,
		Commit:  Commit,
	}

	info, ok := debug.ReadBuildInfo()
	if !ok {
		return result
	}

	result.GoVersion = info.GoVersion

	// Fallback to Go module version.
	// Do not override a version injected by ldflags.
	if result.Version == "dev" &&
		info.Main.Version != "" &&
		info.Main.Version != "(devel)" {
		result.Version = info.Main.Version
	}

	// Fallback to Git commit SHA.
	// Do not override a commit injected by ldflags.
	for _, setting := range info.Settings {
		switch setting.Key {
		case "vcs.revision":
			if result.Commit == "none" {
				result.Commit = setting.Value

				// Short Git SHA.
				if len(result.Commit) > 7 {
					result.Commit = result.Commit[:7]
				}
			}

		case "vcs.modified":
			result.Modified = setting.Value
		}
	}

	return result
}

func main() {
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(versionInfo())
	})

	fmt.Printf("Build info: %+v\n", versionInfo())

	http.ListenAndServe(":8080", nil)
}
