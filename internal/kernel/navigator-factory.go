package kernel

import (
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/pref"
)

type facilities struct {
}

func (f *facilities) Inject(pref.ActiveState) {}

func New(using *pref.Using, o *pref.Options,
	sealer types.GuardianSealer,
	res *types.Resources,
) *Artefacts {
	impl := newImpl(using, o, res)
	controller := newController(using, o, impl, sealer, res)

	return &Artefacts{
		Controller: controller,
		Mediator:   controller.Mediator,
		Resources:  res,
	}
}

func newController(using *pref.Using,
	o *pref.Options,
	impl NavigatorImpl,
	sealer types.GuardianSealer,
	res *types.Resources,
) *NavigationController {
	return &NavigationController{
		Mediator: newMediator(using, o, impl, sealer, res),
	}
}

func newImpl(using *pref.Using,
	o *pref.Options,
	res *types.Resources,
) (impl NavigatorImpl) {
	base := navigator{
		using: using,
		o:     o,
		res:   res,
	}

	switch using.Subscription {
	case enums.SubscribeFiles:
		impl = &navigatorFiles{
			navigator: base,
		}

	case enums.SubscribeFolders, enums.SubscribeFoldersWithFiles:
		impl = &navigatorFolders{
			navigator: base,
		}

	case enums.SubscribeUniversal:
		impl = &navigatorUniversal{
			navigator: base,
		}

	case enums.SubscribeUndefined:
	}

	return
}
