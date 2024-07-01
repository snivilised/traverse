package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
)

func HadesNav(err error) types.KernelController {
	return &navigatorHades{
		err: err,
	}
}

type hadesResult struct {
	err error
}

func (r *hadesResult) Metrics() measure.Reporter {
	return nil
}

func (r *hadesResult) Error() error {
	return r.err
}

type navigatorHades struct {
	err error
}

func (n *navigatorHades) Result(_ context.Context, err error) *types.KernelResult {
	return &types.KernelResult{
		Err: err,
	}
}

func (n *navigatorHades) Starting(types.Session) {
}

func (n *navigatorHades) Navigate(_ context.Context) (core.TraverseResult, error) {
	return &hadesResult{
		err: n.err,
	}, n.err
}
