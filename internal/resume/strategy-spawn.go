package resume

import (
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

func (s *spawnStrategy) resume() (*types.KernelResult, error) {
	return &types.KernelResult{}, nil
}

func (s *spawnStrategy) finish() error {
	return nil
}
