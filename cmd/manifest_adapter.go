package cmd

import "github.com/idzamik/sean/orchestrator"

type ManifestAdapter struct {
	m *Manifest
}

func NewManifestAdapter(m *Manifest) *ManifestAdapter 	{ return &ManifestAdapter{m: m} }

func (a *ManifestAdapter) GetName() string           	{ return a.m.Name }


func (a *ManifestAdapter) GetDefaultCommand() string  	{ return a.m.DefaultCommand }


func (a *ManifestAdapter) GetRules() string           	{ return a.m.Rules }


func (a *ManifestAdapter) GetBin() string 				{ return a.m.Bin }


func (a *ManifestAdapter) GetCommand(name string) (orchestrator.ManifestCommand, bool) {
	cmd, ok := a.m.Commands[name]
	if !ok {
		return orchestrator.ManifestCommand{}, false
	}

	bin := cmd.Bin
	if bin == "" {
		bin = a.m.Bin
	}

	return orchestrator.ManifestCommand{
		Bin:    bin,
		Action: cmd.Action,
		Flags:  cmd.Flags,
	}, true
}


func (a *ManifestAdapter) GetOutput() orchestrator.ManifestOutput {
	return orchestrator.ManifestOutput{
		Format: a.m.Output.Format,
		Dir:    a.m.Output.Dir,
	}
}


func (a *ManifestAdapter) GetCommandNames() []string {
	names := make([]string, 0, len(a.m.Commands))
	for name := range a.m.Commands {
		names = append(names, name)
	}
	return names
}