package kernel

import (
	"context"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/enclave"
)

type navigatorUniversal struct {
	navigatorAgent
}

func (n *navigatorUniversal) Top(ctx context.Context,
	ns *navigationStatic,
) (*enclave.KernelResult, error) {
	return n.top(ctx, ns)
}

func (n *navigatorUniversal) Traverse(ctx context.Context,
	ns *navigationStatic,
	servant core.Servant,
) (bool, error) {
	current := servant.Node()
	descended := ns.mediator.descend(current)

	defer func(permit bool) {
		ns.mediator.ascend(current, permit)
	}(descended)

	if !descended {
		return continueTraversal, nil
	}

	vapour, err := n.inspect(ns, servant)

	if e := ns.mediator.Invoke(servant, vapour); e != nil {
		return continueTraversal, e
	}

	if skip, e := ns.mediator.o.Defects.Skip.Ask(
		current, vapour.Contents(), err,
	); skip == enums.SkipAllTraversal {
		return continueTraversal, e
	} else if skip == enums.SkipDirTraversal {
		return skipTraversal, e
	}

	return n.travel(ctx, ns, vapour)
}

func (n *navigatorUniversal) inspect(ns *navigationStatic,
	servant core.Servant,
) (inspection, error) {
	var (
		current = servant.Node()
		vapour  = &navigationVapour{
			ns:      ns,
			present: current,
		}
		err error
	)

	if current.IsDirectory() {
		vapour.cargo, err = read(ns.mediator.resources.Forest.T,
			n.ro,
			current.Path,
		)

		vapour.Sort(enums.EntryTypeAll)
		vapour.Pick(enums.EntryTypeAll)
	} else {
		vapour.clear()
	}

	extend(ns, vapour)

	return vapour, err
}
