package resume

import (
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type Plugin struct {
	kernel.BasePlugin
	IfResult core.ResultCompletion
}

func (p *Plugin) Next(node *core.Node, inspection types.Inspection) (bool, error) {
	_, _ = node, inspection
	// apply the wake filter

	return true, nil
}

func (p *Plugin) Role() enums.Role {
	return enums.RoleHibernate
}

func (p *Plugin) Init(_ *types.PluginInit) error {
	return p.Mediator.Decorate(p)
}

func (p *Plugin) IsComplete() bool {
	return p.IfResult.IsComplete()
}

func GetSealer(was *pref.Was) types.GuardianSealer {
	if was.Strategy == enums.ResumeStrategyFastward {
		return &fastwardGuardianSealer{}
	}

	return &kernel.Benign{}
}

func Load(res fs.FS, from string, settings ...pref.Option) (*pref.LoadInfo, error) {
	return pref.Load(res, from, settings...)
}
