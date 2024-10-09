package tv

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/measure"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o       *pref.Options
	kc      types.KernelController
	plugins []types.Plugin
	ext     extent
}

// Builders performs build orchestration via its buildAll method. Builders
// is instructed by the factories (via Configure) of which there are 2; one
// for Walk and one for Run. The Prime/Resume extents create the Builders
// instance.
type Builders struct {
	using      *pref.Using
	traverseFS pref.TraverseFileSystemBuilder
	options    optionsBuilder
	navigator  kernel.NavigatorBuilder
	plugins    pluginsBuilder
	extent     extentBuilder
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	// BUILD FILE SYSTEM & EXTENT
	//
	ext := bs.extent.build(
		bs.traverseFS.Build(bs.using.Root),
	)

	// BUILD OPTIONS
	//
	o, binder, optionsErr := bs.options.build(ext)
	if optionsErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, optionsErr),
			ext: ext,
		}, optionsErr
	}

	// BUILD NAVIGATOR
	//
	artefacts, navErr := bs.navigator.Build(o, &types.Resources{
		FS: FileSystems{
			T: ext.traverseFS(),
		},
		Supervisor: measure.New(),
		Binder:     binder,
	})

	if navErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, navErr),
			ext: ext,
		}, navErr
	}

	// BUILD PLUGINS
	//
	plugins, pluginsErr := bs.plugins.build(o,
		bs.using,
		artefacts.Mediator,
		artefacts.Kontroller,
		ext.plugin(artefacts),
	)

	if pluginsErr != nil {
		return &buildArtefacts{
			o:   o,
			kc:  kernel.HadesNav(o, pluginsErr),
			ext: ext,
		}, pluginsErr
	}

	// INIT PLUGINS
	//
	active := lo.Map(plugins,
		func(plugin types.Plugin, _ int) enums.Role {
			return plugin.Role()
		},
	)
	order := manifest(active)
	artefacts.Mediator.Arrange(active, order)
	pi := &types.PluginInit{
		O:        o,
		Controls: &binder.Controls,
	}

	for _, p := range plugins {
		if bindErr := p.Init(pi); bindErr != nil {
			return &buildArtefacts{
				o:       o,
				kc:      artefacts.Kontroller,
				plugins: plugins,
				ext:     ext,
			}, bindErr
		}
	}

	return &buildArtefacts{
		o:       o,
		kc:      artefacts.Kontroller,
		plugins: plugins,
		ext:     ext,
	}, nil
}
