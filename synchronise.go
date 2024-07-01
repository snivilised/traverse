package tv

import (
	"context"

	"github.com/pkg/errors"
	"github.com/snivilised/lorax/boost"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/i18n"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type synchroniser interface {
	core.Navigator
	Starting(types.Session)
	IsComplete() bool
}

type trunk struct {
	nav types.KernelController
	o   *pref.Options
	ext extent
	err error
	u   *pref.Using
	// TODO: w => !!! code-smell: argh, this does not look right (only required for resume)
	w *pref.Was
}

func (t trunk) extent() extent {
	return t.ext
}

func (t trunk) IsComplete() bool {
	return t.ext.complete()
}

func (t trunk) Starting(session types.Session) {
	t.nav.Starting(session)
}

type concurrent struct {
	trunk
	wg        boost.WaitGroup
	pool      *boost.ManifoldFuncPool[*TraverseInput, *TraverseOutput]
	decorator core.Client
	inputCh   boost.SourceStreamW[*TraverseInput]
}

func (c *concurrent) Navigate(ctx context.Context) (core.TraverseResult, error) {
	defer c.close()

	if c.err != nil {
		return c.nav.Result(ctx, c.err), c.err
	}

	c.decorator = func(node *core.Node) error {
		// c.decorator is the function we register with the navigator,
		// so instead of invoking the client handler, the navigator
		// will invoke the decorator, which will send a job to the pool
		// containing the current file system node. The navigator is
		// not aware that its invoking the decorator ...
		//
		// TODO: later, we need to be able to decorate the client callback,
		// either by a Tap or a bus event...
		//
		input := &TraverseInput{
			Node:    node,
			Handler: c.ext.using().Handler,
		}

		c.inputCh <- input // support for timeout (TimeoutOnSendInput) ???

		return nil
	}

	c.pool, c.err = boost.NewManifoldFuncPool(
		ctx, func(input *TraverseInput) (*TraverseOutput, error) {
			err := input.Handler(input.Node)

			return &TraverseOutput{
				Node:  input.Node,
				Error: err,
			}, err
		}, c.wg,
		boost.WithSize(c.o.Core.Concurrency.NoW),
		boost.WithOutput(OutputChSize, CheckCloseInterval, TimeoutOnSend),
	)

	if c.err != nil {
		err := errors.Wrapf(c.err, i18n.ErrWorkerPoolCreationFailed.Error())

		return c.nav.Result(ctx, err), err
	}
	c.open(ctx)

	return c.nav.Navigate(ctx)
}

func (c *concurrent) open(ctx context.Context) {
	c.inputCh = c.pool.Source(ctx, c.wg)
}

func (c *concurrent) close() {
	if c.inputCh != nil {
		close(c.inputCh)
	}
}

type sequential struct {
	trunk
}

func (s *sequential) Navigate(ctx context.Context) (core.TraverseResult, error) {
	if s.err != nil {
		return s.nav.Result(ctx, s.err), s.err
	}

	return s.nav.Navigate(ctx)
}
