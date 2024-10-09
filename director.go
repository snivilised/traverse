package tv

import (
	"github.com/snivilised/traverse/internal/feat/filter"
	"github.com/snivilised/traverse/internal/feat/hiber"
	"github.com/snivilised/traverse/internal/feat/nanny"
	"github.com/snivilised/traverse/internal/feat/resume"
	"github.com/snivilised/traverse/internal/feat/sampling"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/pref"
)

const (
	noOverwrite = true
)

type (
	ifActive func(o *pref.Options,
		using *pref.Using, mediator types.Mediator,
	) types.Plugin
)

// features interrogates options and invokes requests on behalf of the user
// to activate features according to option selections. other plugins will
// be initialised after primary plugins
func features(o *pref.Options, using *pref.Using, mediator types.Mediator,
	kc types.KernelController,
	others ...types.Plugin,
) (plugins []types.Plugin, err error) {
	var (
		all = []ifActive{
			// filtering must happen before sampling so that
			// ReadDirectory hooks are applied to incorrect
			// order. How can we decouple ourselves from this
			// requirement? => the cure is worse than the disease
			//
			hiber.IfActive, nanny.IfActive, filter.IfActive, sampling.IfActive,
		}
	)

	// double reduce, the first reduce 'all' creates list of active plugins
	// and the second, adds other plugins to the activated list.
	plugins = lo.Reduce(others,
		func(acc []types.Plugin, plugin types.Plugin, _ int) []types.Plugin {
			if plugin != nil {
				acc = append(acc, plugin)
			}
			return acc
		},
		lo.Reduce(all,
			func(acc []types.Plugin, query ifActive, _ int) []types.Plugin {
				if plugin := query(o, using, mediator); plugin != nil {
					acc = append(acc, plugin)
				}
				return acc
			},
			[]types.Plugin{},
		),
	)

	for _, plugin := range plugins {
		err = plugin.Register(kc)

		if err != nil {
			return nil, err
		}
	}

	return plugins, nil
}

// Prime extent requests that the navigator performs a full
// traversal from the root path specified.
func Prime(using *pref.Using, settings ...pref.Option) *Builders {
	return &Builders{
		using: using,
		traverseFS: pref.CreateTraverseFS(func(root string) lfs.TraverseFS {
			if using.GetTraverseFS != nil {
				return using.GetTraverseFS(root)
			}

			return lfs.NewTraverseFS(root, noOverwrite)
		}),
		extent: extension(func(tsys lfs.TraverseFS) extent {
			return &primeExtent{
				baseExtent: baseExtent{
					fileSys: fileSystems{
						tsys: tsys,
					},
				},
				u: using,
			}
		}),
		options: optionals(func(ext extent) (*pref.Options, *opts.Binder, error) {
			ve := using.Validate()

			if using.O != nil {
				return using.O, opts.Push(using.O), ve
			}

			o, binder, err := ext.options(settings...)

			if ve != nil {
				return o, binder, ve
			}

			return o, binder, err
		}),
		navigator: kernel.Builder(func(o *pref.Options, // pass in controls here, or put on resources
			resources *types.Resources,
		) (*kernel.Artefacts, error) {
			return kernel.New(using, o, &kernel.Benign{}, resources), nil
		}),
		plugins: activated(features),
	}
}

// Resume extent requests that the navigator performs a resume
// traversal, loading state from a previously saved session
// as a result of it being terminated prematurely via a ctrl-c
// interrupt.
func Resume(was *Was, settings ...pref.Option) *Builders {
	return &Builders{
		using: &was.Using,
		traverseFS: pref.CreateTraverseFS(func(root string) lfs.TraverseFS {
			if was.Using.GetTraverseFS != nil {
				return was.Using.GetTraverseFS(root)
			}

			return lfs.NewTraverseFS(root, noOverwrite)
		}),
		extent: extension(func(tsys lfs.TraverseFS) extent {
			return &resumeExtent{
				baseExtent: baseExtent{
					fileSys: fileSystems{
						tsys: tsys,
					},
				},
				w: was,
			}
		}),
		// we need state; record the hibernation wake point, so
		// using a func here is probably not optimal.
		//
		options: optionals(func(ext extent) (*pref.Options, *opts.Binder, error) {
			ve := was.Validate()

			o, binder, err := ext.options(settings...)

			if ve != nil {
				return o, binder, ve
			}

			return o, binder, err
		}),
		navigator: kernel.Builder(func(o *pref.Options,
			resources *types.Resources,
		) (*kernel.Artefacts, error) {
			return resume.NewController(was,
				kernel.New(&was.Using, o, resume.GetSealer(was), resources),
			), nil
		}),
		plugins: activated(features),
	}
}
