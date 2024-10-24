package nanny

// 📦 pkg: nanny - handles a node's children for folders with children subscription

import (
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

func IfActive(o *pref.Options,
	using *pref.Using, mediator types.Mediator,
) types.Plugin {
	if using.Subscription == enums.SubscribeFoldersWithFiles &&
		!o.Filter.IsFilteringActive() {
		return &plugin{
			BasePlugin: kernel.BasePlugin{
				O:             o,
				Mediator:      mediator,
				ActivatedRole: enums.RoleNanny,
			},
		}
	}

	return nil
}

type plugin struct {
	kernel.BasePlugin
	crate measure.Crate
}

func (p *plugin) Next(servant core.Servant,
	inspection types.Inspection,
) (bool, error) {
	node := servant.Node()
	files := inspection.Sort(enums.EntryTypeFile)
	node.Children = files
	p.crate.Mums[enums.MetricNoChildFilesFound].Times(uint(len(files)))

	return true, nil
}

func (p *plugin) Init(_ *types.PluginInit) error {
	p.crate.Mums = p.Mediator.Supervisor().Many(
		enums.MetricNoChildFilesFound,
	)

	return p.Mediator.Decorate(p)
}
