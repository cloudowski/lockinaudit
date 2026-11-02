// Package version holds the versions that identify every lockinaudit result.
package version

import (
	"runtime/debug"
	"strings"
)

// devVersion is the Tool value of a build that was given no version.
const devVersion = "0.2.0-dev"

// Tool is the lockinaudit version, set at build time via
// -ldflags "-X github.com/cloudowski/lockinaudit/internal/version.Tool=..."
// Builds without that flag, such as `go install ...@v0.2.0`, use the module version recorded
// by the Go toolchain (see moduleVersion), and a plain `go build` keeps the development version.
var Tool = devVersion

func init() {
	if Tool != devVersion {
		return
	}
	if info, ok := debug.ReadBuildInfo(); ok {
		if v := moduleVersion(info); v != "" {
			Tool = v
		}
	}
}

// moduleVersion returns the version of the main module as a release version without the
// leading "v" (the form GoReleaser stamps, for example "0.2.0"), or "" when the build has none:
// a build from a source checkout reports "(devel)".
func moduleVersion(info *debug.BuildInfo) string {
	if info == nil {
		return ""
	}
	v := info.Main.Version
	if v == "" || v == "(devel)" {
		return ""
	}
	return strings.TrimPrefix(v, "v")
}

// SupportEmail is shown at the bottom of reports. It is a build-time setting:
// -ldflags "-X github.com/cloudowski/lockinaudit/internal/version.SupportEmail=..."
var SupportEmail = ""

// RulesVersion is the version of the rule evaluation and scoring semantics.
// Bump it whenever rule evaluation or any scoring formula changes (research R19).
const RulesVersion = "1.1.0"

// ResultFormat identifies the result document format.
const ResultFormat = "lockinaudit.result/v1"

// SupportedCatalogSchemaMajors lists catalog schema majors this build can read.
var SupportedCatalogSchemaMajors = []int{1}
