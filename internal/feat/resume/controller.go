package resume

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
	"github.com/snivilised/agenor/internal/kernel"
	"github.com/snivilised/agenor/internal/opts"
	"github.com/snivilised/agenor/pref"
)

type Controller struct {
	med      enclave.Mediator
	relic    *pref.Relic
	load     *opts.LoadInfo
	strategy Strategy
}

func (c *Controller) Ignite(ignition *enclave.Ignition) {
	c.med.Ignite(ignition)
}

func (c *Controller) Result(ctx context.Context, err error) *enclave.KernelResult {
	return c.med.Result(ctx, err)
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
	c.med.Conclude(result)
}

func newStrategy(inception *kernel.Inception,
	sealer enclave.GuardianSealer,
	mediator enclave.Mediator,
) (strategy Strategy) {
	load := inception.Harvest.Loaded()
	relic, _ := inception.Facade.(*pref.Relic)
	base := baseStrategy{
		o:        load.O,
		active:   load.State,
		relic:    relic,
		sealer:   sealer,
		kc:       mediator,
		mediator: mediator,
		forest:   inception.Resources.Forest,
	}

	switch relic.Strategy {
	case enums.ResumeStrategyFastward, enums.ResumeStrategyUndefined:
		strategy = &fastwardStrategy{
			baseStrategy: base,
			role:         enums.RoleFastward,
		}

	case enums.ResumeStrategySpawn:
		strategy = &spawnStrategy{
			baseStrategy: base,
		}
	}

	return strategy
}

func (c *Controller) Navigate(ctx context.Context) (*enclave.KernelResult, error) {
	if err := c.strategy.init(c.load); err != nil {
		return c.Result(ctx, err), err
	}

	return c.strategy.resume(ctx)
}
