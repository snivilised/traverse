package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/pref"
)

type Controller struct {
	kc         enclave.KernelController
	was        *pref.Was
	load       *opts.LoadInfo
	strategy   Strategy
	facilities enclave.Facilities
}

func (c *Controller) Ignite(ignition *enclave.Ignition) {
	c.kc.Ignite(ignition)
}

func (c *Controller) Result(ctx context.Context, err error) *enclave.KernelResult {
	return c.kc.Result(ctx, err)
}

func (c *Controller) Mediator() enclave.Mediator {
	return c.kc.Mediator()
}

func (c *Controller) Strategy() Strategy {
	return c.strategy
}

func (c *Controller) Resume(context.Context,
	*core.ActiveState,
) (*enclave.KernelResult, error) {
	return &enclave.KernelResult{}, nil
}

func (c *Controller) Conclude(result core.TraverseResult) {
	c.kc.Conclude(result)
}

func newStrategy(was *pref.Was,
	harvest enclave.OptionHarvest,
	kc enclave.KernelController,
	sealer enclave.GuardianSealer,
	resources *enclave.Resources,
) (strategy Strategy) {
	load := harvest.Loaded()
	base := baseStrategy{
		o:        load.O,
		active:   load.State,
		was:      was,
		sealer:   sealer,
		kc:       kc,
		mediator: kc.Mediator(),
		forest:   resources.Forest,
	}

	switch was.Strategy {
	case enums.ResumeStrategyFastward:
		strategy = &fastwardStrategy{
			baseStrategy: base,
			role:         enums.RoleFastward,
		}

	case enums.ResumeStrategySpawn:
		strategy = &spawnStrategy{
			baseStrategy: base,
		}
	case enums.ResumeStrategyUndefined:
	}

	return strategy
}

func (c *Controller) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	if err := c.strategy.init(c.load); err != nil {
		return c.Result(ctx, err), err
	}

	return c.strategy.resume(ctx, c.was)
}
