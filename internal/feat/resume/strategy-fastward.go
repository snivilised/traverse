package resume

import (
	"context"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
)

type fastwardGuardianSealer struct {
}

func (g *fastwardGuardianSealer) Seal(top types.Link) error {
	if top.Role() == enums.RoleHibernate {
		return core.ErrGuardianCantDecorateItemSealed
	}
	return nil
}

func (g *fastwardGuardianSealer) IsSealed(top types.Link) bool {
	_ = top

	return false
}

type fastwardStrategy struct {
	baseStrategy
}

func (s *fastwardStrategy) init() {

}

func (s *fastwardStrategy) attach() {

}

func (s *fastwardStrategy) detach() {

}

func (s *fastwardStrategy) resume(context.Context) (*types.KernelResult, error) {
	return &types.KernelResult{}, nil
}

func (s *fastwardStrategy) ifResult() bool {
	return true
}

func (s *fastwardStrategy) finish() error {
	return nil
}
