package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/level"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/pref"
)

// mediator controls traversal events, sends notifications and emits
// life-cycle events
type mediator struct {
	root         string
	subscription enums.Subscription
	using        *pref.Using
	impl         NavigatorImpl
	guardian     *guardian
	periscope    *level.Periscope
	o            *pref.Options
	resources    *types.Resources
	mums         measure.MutableMetrics
	order        []enums.Role
}

func newMediator(using *pref.Using,
	o *pref.Options,
	impl NavigatorImpl,
	sealer types.GuardianSealer,
	resources *types.Resources,
) *mediator {
	mums := resources.Supervisor.Many(
		enums.MetricNoFilesInvoked,
		enums.MetricNoFoldersInvoked,
		enums.MetricNoChildFilesFound,
	)

	return &mediator{
		root:         using.Root,
		subscription: using.Subscription,
		using:        using,
		impl:         impl,
		guardian:     newGuardian(using.Handler, sealer, mums),
		periscope:    level.New(),
		o:            o,
		resources:    resources,
		mums:         mums,
	}
}

func (m *mediator) descend(node *core.Node) bool {
	if !m.periscope.Descend(m.o.Behaviours.Cascade.Depth) {
		return false
	}

	m.o.Binder.Controls.Descend.Dispatch()(node)

	return true
}

func (m *mediator) ascend(node *core.Node, permit bool) {
	if permit {
		m.periscope.Ascend()
		m.o.Binder.Controls.Ascend.Dispatch()(node)
	}
}

func (m *mediator) Decorate(link types.Link) error {
	return m.guardian.Decorate(link)
}

func (m *mediator) Unwind(role enums.Role) error {
	return m.guardian.Unwind(role)
}

func (m *mediator) Arrange(active, order []enums.Role) {
	m.order = order
	m.guardian.arrange(active, order)
}

func (m *mediator) Ignite(ignition *types.Ignition) {
	m.impl.Ignite(ignition)
	m.o.Binder.Controls.Begin.Dispatch()(&cycle.BeginState{
		Root: m.root,
	})
}

func (m *mediator) Conclude(result core.TraverseResult) {
	m.o.Binder.Controls.End.Dispatch()(result)
}

func (m *mediator) Navigate(ctx context.Context) (core.TraverseResult, error) {
	result, err := m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		root:     m.root,
	})

	return result, err
}

func (m *mediator) Spawn(ctx context.Context, root string) (core.TraverseResult, error) {
	// TODO: send a message indicating spawn
	//
	return m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		root:     root,
	})
}

func (m *mediator) Invoke(node *core.Node, inspection core.Inspection) error {
	return m.guardian.Invoke(node, inspection)
}

func (m *mediator) Supervisor() *measure.Supervisor {
	return m.resources.Supervisor
}
