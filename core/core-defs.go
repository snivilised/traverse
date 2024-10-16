package core

import (
	"time"
)

// 📦 pkg: core - contains universal definitions and handles user facing cross
// cutting concerns try to keep to a minimum to reduce rippling changes.

type (
	// ResultCompletion used to determine if the result really represents
	// final navigation completion.
	ResultCompletion interface {
		IsComplete() bool
	}

	Completion func() bool

	// Session represents a traversal session and keeps tracks of
	// timing.
	Session interface {
		ResultCompletion
		StartedAt() time.Time
		Elapsed() time.Duration
	}

	// TraverseResult
	TraverseResult interface {
		Metrics() Reporter
		Session() Session
		Error() error
	}

	// Servant provides the client with facility to request properties
	// about the current navigation node.
	Servant interface {
		Node() *Node
	}

	// Client is the callback invoked for each file system node found
	// during traversal.
	Client func(servant Servant) error

	// SimpleHandler is a function that takes no parameters and can
	// be used by any notification with this signature.
	SimpleHandler func()

	// BeginHandler invoked before traversal begins
	BeginHandler func(tree string)

	// EndHandler invoked at the end of traversal
	EndHandler func(result TraverseResult)

	// HibernateHandler is a generic handler that is used by hibernation
	// to indicate wake or sleep.
	HibernateHandler func(description string)
)

func (fn Completion) IsComplete() bool {
	return fn()
}
