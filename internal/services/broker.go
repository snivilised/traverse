package services

import "github.com/snivilised/extendio/bus"

const (
	format                  = "%03d"
	TopicInitPlugins        = "topic:init.plugins"
	TopicInterceptNavigator = "topic:intercept.navigator"
	TopicNavigationComplete = "topic:navigation.complete"
	TopicOptionsAnnounce    = "topic:options.announce"
	TopicOptionsBefore      = "topic:options.before"
	TopicOptionsComplete    = "topic:options.complete"
)

var (
	Broker *bus.Broker
	topics = []string{
		TopicInitPlugins,
		TopicInterceptNavigator,
		TopicNavigationComplete,
		TopicOptionsAnnounce,
		TopicOptionsBefore,
		TopicOptionsComplete,
	}
)

func init() {
	Reset()
}

type (
	InitBroker interface {
		Available(b *bus.Broker)
	}
)

func Reset() *bus.Broker {
	b, err := bus.New(&bus.Sequential{
		Format: format,
	})

	if err != nil {
		panic(err)
	}

	b.RegisterTopics(topics...)

	// Access to the broker is currently un-synchronised; that is because interaction
	// with the broker is only expected to come from a single thread. However, if we
	// really wanted to grant access to it to other threads, we can define a wrapper
	// function/object around it that implements synchronisation using locks.
	//
	Broker = b

	return b
}
