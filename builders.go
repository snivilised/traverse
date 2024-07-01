package tv

import (
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/pref"
)

type buildArtefacts struct {
	o       *pref.Options
	nav     types.KernelController
	plugins []types.Plugin
	ext     extent
}

type Builders struct {
	filesystem pref.FileSystemBuilder
	options    optionsBuilder
	navigator  kernel.NavigatorBuilder
	plugins    pluginsBuilder
	extent     extentBuilder
}

func (bs *Builders) buildAll() (*buildArtefacts, error) {
	// BUILD FILE SYSTEM & EXTENT
	//
	ext := bs.extent.build(bs.filesystem.Build())

	// BUILD OPTIONS
	//
	o, optionsErr := bs.options.build(ext)
	if optionsErr != nil {
		return &buildArtefacts{
			o:   o,
			nav: kernel.HadesNav(optionsErr),
			ext: ext,
		}, optionsErr
	}

	// BUILD NAVIGATOR
	//
	artefacts, navErr := bs.navigator.Build(o, &types.Resources{
		FS: types.FileSystems{
			N: ext.navFS(),
			R: ext.resFS(),
		},
		Supervisor: measure.New(),
	})
	if navErr != nil {
		return &buildArtefacts{
			o:   o,
			nav: kernel.HadesNav(navErr),
			ext: ext,
		}, navErr
	}

	// BUILD PLUGINS
	//
	plugins, pluginsErr := bs.plugins.build(o,
		artefacts.Mediator,
		ext.plugin(artefacts.Mediator),
	)

	if pluginsErr != nil {
		return &buildArtefacts{
			o:   o,
			nav: kernel.HadesNav(pluginsErr),
			ext: ext,
		}, pluginsErr
	}

	// INIT PLUGINS
	//
	for _, p := range plugins {
		if bindErr := p.Init(); bindErr != nil {
			return &buildArtefacts{
				o:       o,
				nav:     artefacts.Controller,
				plugins: plugins,
				ext:     ext,
			}, bindErr
		}
	}

	return &buildArtefacts{
		o:       o,
		nav:     artefacts.Controller,
		plugins: plugins,
		ext:     ext,
	}, nil
}
