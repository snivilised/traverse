package tv

import (
	"context"

	"github.com/snivilised/extendio/bus"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
)

const (
	badge = "badge: navigation-driver"
)

type driver struct {
	s session
}

func (d *driver) init() {
	services.Broker.RegisterHandler(badge, bus.Handler{
		Handle: func(_ context.Context, m bus.Message) {
			_ = m.Data
			// now invoke session.finish
		},
		Matcher: services.TopicNavigationComplete,
	})
}

func (d *driver) Navigate(ctx context.Context) (core.TraverseResult, error) {
	d.init()
	d.s.start()
	result, err := d.s.exec(ctx)

	d.s.finish(result)

	return result, err
}
