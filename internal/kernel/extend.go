package kernel

import (
	"path/filepath"
	"strings"

	"github.com/snivilised/agenor/core"
	"github.com/snivilised/agenor/enums"
	"github.com/snivilised/agenor/internal/third/lo"
)

func extend(ns *navigationStatic, vapour inspection) {
	var (
		scope    enums.FilterScope
		isLeaf   bool
		current  = vapour.Current()
		contents = vapour.Contents()
	)

	if current.IsDirectory() {
		isLeaf = len(contents.Directories()) == 0
		scope = ns.mediator.periscope.Scope(isLeaf)
		scope |= enums.ScopeDirectory
	} else {
		scope = enums.ScopeLeaf
		scope |= enums.ScopeFile
	}

	parent, name := filepath.Split(current.Path)
	current.Extension = core.Extension{
		Depth:  ns.mediator.periscope.Depth(),
		IsLeaf: isLeaf,
		Name:   name,
		Parent: parent,
		Scope:  scope,
	}

	keepTrailingSep := ns.mediator.o.Behaviours.SubPath.KeepTrailingSep

	spInfo := &core.SubPathInfo{
		Tree:            ns.tree,
		Node:            current,
		KeepTrailingSep: keepTrailingSep,
	}

	subpath := lo.TernaryF(current.IsDirectory(),
		func() string { return ns.mediator.o.Hooks.DirectorySubPath.Invoke()(spInfo) },
		func() string { return ns.mediator.o.Hooks.FileSubPath.Invoke()(spInfo) },
	)

	subpath = lo.TernaryF(keepTrailingSep,
		func() string { return subpath },
		func() string {
			result := subpath
			sep := string(filepath.Separator)

			if strings.HasSuffix(subpath, sep) {
				result = subpath[:strings.LastIndex(subpath, sep)]
			}

			return result
		},
	)

	current.Extension.SubPath = subpath
}
