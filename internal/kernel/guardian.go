package kernel

import (
	"errors"

	"github.com/snivilised/extendio/collections"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/measure"
)

type (
	invocationChain = collections.Stack[types.Link]
	invocationIt    = collections.Iterator[types.Link]
)

type owned struct {
	mums measure.Mutables
}

// anchor is a specialised link that should always be the
// last in the chain and contains the original client's handler.
type anchor struct {
	target core.Client
	owned  owned
}

func (t *anchor) Next(node *core.Node) (bool, error) {
	if metric := lo.Ternary(node.IsFolder(),
		t.owned.mums[enums.MetricNoFoldersInvoked],
		t.owned.mums[enums.MetricNoFilesInvoked],
	); metric != nil {
		metric.Tick()
	}

	return false, t.target(node)
}

func (t *anchor) Role() enums.Role {
	return enums.RoleTerminus
}

type iterationContainer struct {
	invocationChain
	// it iterates the contents of the chain in reverse; we
	// need a reverse iterator, because the first element of the chain
	// (index 0), will always represent the client's actual callback,
	// that should only ever be called last, if the prior links permits it
	// to be so.
	//
	it invocationIt
}

// guardian controls access to the client callback
type guardian struct {
	callback core.Client
	chain    iterationContainer
	master   types.GuardianSealer
}

func newGuardian(callback core.Client,
	master types.GuardianSealer,
	mums measure.Mutables,
) *guardian {
	// TODO: need to pass in a sequence manifest that describes the
	// valid chain of decorators, defined by role eg
	// { enums.RoleTerminus, enums.RoleClientFilter, }
	//
	stack := collections.NewStack[types.Link]()
	stack.Push(&anchor{
		target: callback,
		owned: owned{
			mums: mums,
		},
	})

	return &guardian{
		callback: callback,
		chain: iterationContainer{
			invocationChain: *stack,
			it:              collections.ReverseIt(stack.Content(), nil),
		},
		master: master,
	}
}

// role indicates the guise under which the decorator is being applied.
// Not all roles can be decorated (sealed). The fastward-resume decorator is
// sealed. If an attempt is made to Decorate a sealed decorator,
// an error is returned.
func (g *guardian) Decorate(link types.Link) error {
	// role enums.Role, decorator core.Client
	//
	// if every feature is active
	// [prime]:
	// hiber => sampling => filtering => callback
	// [fastward-resume]:
	// fastward-filter => [prime]
	//
	// sequence: fastward-filter => hiber => sampling => filtering => callback
	//
	top, err := g.chain.Current()
	if err != nil {
		return err
	}

	if g.master.IsSealed(top) {
		return errors.New("can't decorate, last item is sealed")
	}

	g.chain.Push(link)
	g.create()

	return nil
}

func (g *guardian) Unwind(enums.Role) error {
	return nil
}

// Invoke executes the chain which may or may not end up resulting in
// the invocation of the client's callback, depending on the contents
// of the chain.
func (g *guardian) Invoke(node *core.Node) error {
	// TODO: Actually, using an iterator is not the best way forward as it
	// adds unnecessary overhead. Each link should have access to the next,
	// without depending on an iterator.
	for link := g.chain.it.Start(); g.chain.it.Valid(); g.chain.it.Next() {
		if next, err := link.Next(node); !next || err != nil {
			return err
		}
	}

	return nil
}

// create rebuilds the iterator, should be called after the chain has been
// modified.
func (g *guardian) create() {
	g.chain.it = collections.ReverseIt(g.chain.Content(), nil)
}

// Benign is used when a master sealer has not been registered. It is
// permissive in nature.
type Benign struct {
}

func (m *Benign) Seal(types.Link) error {
	return nil
}

func (m *Benign) IsSealed(types.Link) bool {
	return false
}
