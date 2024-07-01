package resume

import (
	"context"

	"github.com/snivilised/traverse/internal/kernel"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

const (
	badge = "badge: resume"
)

// refine should also contain persistence concerns (actually
// these may be internal modules, eg internal/serial/JSON). Depends on hiber, refine
// and persist.

type resumeStrategy interface {
	init()
	attach()
	detach()
	resume(context.Context) (*types.KernelResult, error)
	complete() bool
	finish() error
}

type baseStrategy struct {
	o    *pref.Options
	kc   types.KernelController
	impl kernel.NavigatorImpl
}
