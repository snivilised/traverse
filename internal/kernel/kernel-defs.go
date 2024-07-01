package kernel

import (
	"context"
	"io/fs"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type (
	// NavigatorImpl
	NavigatorImpl interface {
		Starting(session types.Session)

		// Top
		Top(ctx context.Context,
			ns *navigationStatic,
		) (*types.KernelResult, error)

		// Travel is the internal version of Traverse. It is useful to distinguish
		// between the external Traverse and the internal Travel, because they
		// have different return types and semantics.
		Travel(ctx context.Context,
			ns *navigationStatic,
			current *core.Node,
		) (bool, error)

		Result(ctx context.Context, err error) *types.KernelResult
	}

	// NavigatorDriver
	NavigatorDriver interface {
		Impl() NavigatorImpl
	}

	// Gateway provides a barrier around the guardian to prevent accidental
	// misuse.
	Gateway interface {
		// Decorate is used to wrap the existing client. The decoration will
		// result in the decorator being called first before the existing the
		// client. This needs to behave like chain of responsibility, because
		// the decorator has the choice of wether to pass on the call further
		// down the chain and ultimately the client's callback.
		// The returned bool indicates wether the next in the chain is invoked
		// or not; true, means pass on, false means absorb.
		//
		// A filter decorator will return true if the node matches filter, false
		// otherwise.
		// A fastward resume will act likewise. If we have fastward active, but
		// there is also an underlying filter we have a chain that looks like
		// this:
		// fastward-filter => underlying-filter => callback
		//
		// With this in place, we have to think very carefully as to whether we
		// really need a state machine, because the chain is able to fulfill this
		// purpose.
		//
		// but wait, let's think about wake and sleep. In the normal scenario,
		// fastward will start off in sleeping mode and the filter will be a
		// wake condition. Once we encounter the wake condition, the hibernation
		// decorator needs to be removed. But how do we know that the top of the
		// chain is the hibernate decorator? It must be that the resume hibernate
		// can't be decorated, this suggests some kind of priority/authorisation
		// is required. Because of this, we need a role and we also need to make
		// sure the features are initialised in the correct order to make this
		// happen correctly => sequence manifest
		//
		// role indicates the guise under which the decorator is being applied.
		// Not all roles can be decorated. The fastward-resume decorator can
		// not be decorated. If an attempt is made to Decorate a sealed decorator,
		// an error is returned.
		Decorate(link types.Link) error

		// Invoke executes the chain which may or may not end up resulting in
		// the invocation of the client's callback, depending on the contents
		// of the chain.
		// Invoke(node *core.Node) error

		// Unwind removes last link in the chain which is expected to be of
		// role specified.
		Unwind(role enums.Role) error
	}

	Invokable interface {
		Gateway
		Invoke(node *core.Node) error
	}

	// navigationStatic contains static info, ie info that is established during
	// bootstrap and doesn't change after navigation begins. Used to help
	// minimise allocations.
	navigationStatic struct {
		mediator *mediator
		root     string
	}

	// navigationVapour represents short-lived navigation data whose state relates
	// only to the current Node. (equivalent to inspection in extendio)
	navigationVapour struct { // after content has been read
		ns                *navigationStatic
		currentNode       *core.Node
		directoryContents *DirectoryContents
		ents              []fs.DirEntry
	}

	navigationInfo struct { // pre content read
	}

	inspection interface { // after content has been read
		static() *navigationStatic
		current() *core.Node
		contents() core.DirectoryContents
		entries() []fs.DirEntry
		clear()
	}

	navigationAssets struct {
		ns     navigationStatic
		vapour *navigationVapour
	}
)

func (v *navigationVapour) static() *navigationStatic {
	return v.ns
}

func (v *navigationVapour) current() *core.Node {
	return v.currentNode
}

func (v *navigationVapour) contents() core.DirectoryContents {
	return v.directoryContents
}

func (v *navigationVapour) entries() []fs.DirEntry {
	return v.ents
}

func (v *navigationVapour) clear() {
	if v.directoryContents != nil {
		v.directoryContents.Clear()
	} else {
		newEmptyDirectoryEntries(v.ns.mediator.o)
	}
}
