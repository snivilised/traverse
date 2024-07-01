package kernel

import (
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type (
	Artefacts struct {
		Controller types.KernelController
		Mediator   types.Mediator
		Facilities types.Facilities
		Resources  *types.Resources
	}

	NavigatorBuilder interface {
		Build(o *pref.Options, res *types.Resources) (*Artefacts, error)
	}

	Builder func(o *pref.Options, res *types.Resources) (*Artefacts, error)
)

func (fn Builder) Build(o *pref.Options, res *types.Resources) (*Artefacts, error) {
	return fn(o, res)
}
