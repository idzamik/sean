
package meta

const (
	AppName = "sean"
	AppFullName = "Security Analyser"
	AppVersion = "0.1.0"
	AppDescription = AppName + " — " + AppFullName + ": unified CLI orchestrator for SAST, SCA and other security analysis tools"

	AppLongDescription = AppName + ` (` + AppFullName + `) orchestrates security analysis tools
in offline/air-gap environments. Tools are deployed from local distributives
using declarative YAML manifests — no internet access required.

Use "` + AppName + ` setup install <tool>" to deploy a tool,
then "` + AppName + ` sast" or "` + AppName + ` sca" to run analysis.`
)
