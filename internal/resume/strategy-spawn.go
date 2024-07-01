package resume

import (
	"context"

	"github.com/snivilised/traverse/internal/types"
)

type spawnStrategy struct {
	baseStrategy
}

func (s *spawnStrategy) init() {

}

func (s *spawnStrategy) attach() {

}

func (s *spawnStrategy) detach() {

}

func (s *spawnStrategy) resume(ctx context.Context) (*types.KernelResult, error) {
	return s.impl.Result(ctx, nil), nil
}

func (s *spawnStrategy) finish() error {
	return nil
}

func (s *spawnStrategy) complete() bool {
	panic("NOT-IMPL:spawnStrategy.complete")
}
