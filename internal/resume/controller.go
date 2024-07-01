package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/i18n"
	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type Controller struct {
	kc         types.KernelController
	was        *pref.Was
	strategy   resumeStrategy
	facilities types.Facilities
}

func (c *Controller) Starting(session types.Session) {
	c.kc.Starting(session)
}

func (c *Controller) Result(ctx context.Context, err error) *types.KernelResult {
	return c.kc.Result(ctx, err)
}

func NewController(was *pref.Was, artefacts *kernel.Artefacts) *kernel.Artefacts {
	// The Navigator on the incoming artefacts is the core navigator. It is
	// decorated here for resume. The strategy only needs access to the core navigator.
	// The resume navigator delegates to the strategy.
	//
	var (
		strategy resumeStrategy
		err      error
	)

	if strategy, err = newStrategy(was, artefacts.Controller); err != nil {
		return artefacts
	}

	return &kernel.Artefacts{
		Controller: &Controller{
			kc:         artefacts.Controller,
			was:        was,
			strategy:   strategy,
			facilities: artefacts.Facilities,
		},
		Mediator: artefacts.Mediator,
	}
}

func newStrategy(was *pref.Was, kc types.KernelController) (strategy resumeStrategy, err error) {
	driver, ok := kc.(kernel.NavigatorDriver)

	if !ok {
		return nil, i18n.ErrInternalFailedToGetNavigatorDriver
	}

	base := baseStrategy{
		o:    was.O,
		kc:   kc,
		impl: driver.Impl(),
	}

	switch was.Strategy {
	case enums.ResumeStrategyFastward:
		strategy = &fastwardStrategy{
			baseStrategy: base,
		}
	case enums.ResumeStrategySpawn:
		strategy = &spawnStrategy{
			baseStrategy: base,
		}
	case enums.ResumeStrategyUndefined:
	}

	return strategy, nil
}

func (c *Controller) Navigate(ctx context.Context) (core.TraverseResult, error) {
	return c.strategy.resume(ctx)
}
