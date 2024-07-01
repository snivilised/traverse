package kernel

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/level"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
	"github.com/snivilised/traverse/pref"
)

// mediator controls traversal events, sends notifications and emits life-cycle events

type mediator struct {
	root      string
	impl      NavigatorImpl
	client    Invokable
	frame     *navigationFrame
	pad       *scratchPad // gets created just before nav begins
	o         *pref.Options
	resources *types.Resources
	// there should be a registration phase; but doing so mean that
	// these entities should already exist, which is counter productive.
	// possibly use dependency inject where entities declare their
	// dependencies so that it is easier to orchestrate the boot.
	//
	// Are hooks plugins? (
	// query-status: lstat on the root directory; singular; WithQueryStatus
	// read-dir: ; singular; WithReader
	// folder-sub-path: ; singular; WithFolderSubPath
	// filter-int: ; singular; WithFilter
	// sort: ; singular
	// )
	//
	// Hooks mark specific points in navigation that can be customised
}

func newMediator(using *pref.Using,
	o *pref.Options,
	impl NavigatorImpl,
	sealer types.GuardianSealer,
	res *types.Resources,
) *mediator {
	return &mediator{
		root: using.Root,
		impl: impl,
		client: newGuardian(using.Handler, sealer, res.Supervisor.Many(
			enums.MetricNoFilesInvoked,
			enums.MetricNoFoldersInvoked,
		)),
		frame: &navigationFrame{
			periscope: level.New(),
		},
		pad:       newScratch(o),
		o:         o,
		resources: res,
	}
}

func (m *mediator) Decorate(link types.Link) error {
	return m.client.Decorate(link)
}

func (m *mediator) Unwind(role enums.Role) error {
	return m.client.Unwind(role)
}

func (m *mediator) Starting(session types.Session) {
	m.impl.Starting(session)
}

func (m *mediator) Navigate(ctx context.Context) (core.TraverseResult, error) {
	// could we pass in the invokable client to Top so the navigators can invoke
	// as required.
	//
	result, err := m.impl.Top(ctx, &navigationStatic{
		mediator: m,
		root:     m.root,
	})

	result.Reporter = m.resources.Supervisor

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

func (m *mediator) Invoke(node *core.Node) error {
	return m.client.Invoke(node)
}

func (m *mediator) Supervisor() *measure.Supervisor {
	return m.resources.Supervisor
}

// application phases (should we define a state machine?)
//
// --> configuration: OnConfigured
// * ...
// --> i18n: Oni18n
// * ...
// --> log
// * ...
// --> session
// *
// --> get-options (via WithOptions for primary session, or restore with resume session)
