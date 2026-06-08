// Package meta содержит все идентификаторы приложения.
// Чтобы переименовать CLI — достаточно изменить константы только в этом файле.
// Все команды, сообщения и help-тексты подтягивают имя отсюда.
package meta

const (
	// AppName — отображаемое имя CLI во всех сообщениях и help-тексте.
	AppName = "sean"

	// AppFullName — расшифровка аббревиатуры.
	AppFullName = "Security Analyser"

	// AppVersion — текущая версия приложения.
	AppVersion = "0.1.0"

	// AppDescription — краткое описание для корневой команды Cobra.
	AppDescription = AppName + " — " + AppFullName + ": unified CLI orchestrator for SAST, SCA and other security analysis tools"

	// AppLongDescription — полное описание для --help.
	AppLongDescription = AppName + ` (` + AppFullName + `) orchestrates security analysis tools
in offline/air-gap environments. Tools are deployed from local distributives
using declarative YAML manifests — no internet access required.

Use "` + AppName + ` setup install <tool>" to deploy a tool,
then "` + AppName + ` sast" or "` + AppName + ` sca" to run analysis.`
)
