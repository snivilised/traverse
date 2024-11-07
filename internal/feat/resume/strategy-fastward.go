package resume

import (
	"context"
	"fmt"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/enclave"
	"github.com/snivilised/traverse/internal/filtering"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/pref"
)

type FastwardFilter struct {
	source      string
	description string
	scope       enums.FilterScope
	parent      string
	name        string
}

// Description describes filter
func (f *FastwardFilter) Description() string {
	return f.description
}

// Validate ensures the filter definition is valid, returns
// error when invalid
func (f *FastwardFilter) Validate() error {
	return nil
}

// Source, filter definition (comes from filter definition Pattern)
func (f *FastwardFilter) Source() string {
	return f.source
}

// IsMatch does this node match the filter
func (f *FastwardFilter) IsMatch(node *core.Node) bool {
	return node.Extension.Name == f.name && node.Parent.Path == f.parent
}

// IsApplicable is this filter applicable to this node's scope
func (f *FastwardFilter) IsApplicable(node *core.Node) bool {
	return (node.IsDirectory() && f.scope.IsDirectory()) ||
		(!node.IsDirectory() && f.scope.IsFile())
}

// Scope, what items this filter applies to
func (f *FastwardFilter) Scope() enums.FilterScope {
	return f.scope
}

type fastwardGuardianSealer struct {
}

func (g *fastwardGuardianSealer) Seal(top enclave.Link) error {
	if top.Role() == enums.RoleHibernate {
		return core.ErrGuardianCantDecorateItemSealed
	}

	return nil
}

func (g *fastwardGuardianSealer) IsSealed(enclave.Link) bool {
	return false
}

type fastwardStrategy struct {
	baseStrategy
	enclave.Link
	role   enums.Role
	filter core.TraverseFilter
}

func (s *fastwardStrategy) init(load *opts.LoadInfo) (err error) {
	// We don't use the Hibernate.Wake/Sleep-At, as those are defined bt the client.
	// Instead we just need to create a filter on the fly from the load state...
	//
	calc := s.forest.T.Calc()
	scope := lo.Ternary(load.State.IsDir, enums.ScopeDirectory, enums.ScopeFile)
	s.filter, err = filtering.New(
		&core.FilterDef{
			Type: enums.FilterTypeCustom,
		},
		&pref.FilterOptions{
			Custom: &FastwardFilter{
				description: fmt.Sprintf("[scope: '%v'], path: '%v'",
					scope, load.State.CurrentPath,
				),
				scope:  scope,
				source: load.State.CurrentPath,
				parent: calc.Dir(load.State.CurrentPath),
				name:   calc.Base(load.State.CurrentPath),
			},
		},
	)

	if err := s.mediator.Decorate(s); err != nil {
		return err
	}

	return err
}

// Next invokes this decorator which returns true if
// next link in the chain can be run or false to stop
// execution of subsequent links.
func (s *fastwardStrategy) Next(servant core.Servant,
	_ enclave.Inspection,
) (match bool, err error) {
	node := servant.Node()
	match = s.filter.IsMatch(node)

	if match {
		// TODO: unmute notifications
		//
		err = s.mediator.Unwind(s.role)
	}

	return match, err
}

// Role indicates the identity of the link
func (s *fastwardStrategy) Role() enums.Role {
	return enums.RoleFastward
}

func (s *fastwardStrategy) resume(ctx context.Context) (*enclave.KernelResult, error) {
	return s.mediator.Resume(ctx, s.active)
}

func (s *fastwardStrategy) ifResult() bool {
	return true
}
