package refine

import (
	"slices"
	"strings"

	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

func fromExtendedGlobPattern(pattern string) (segments, suffixes []string, err error) {
	if !strings.Contains(pattern, "|") {
		return []string{}, []string{},
			locale.NewInvalidExtGlobFilterMissingSeparatorError(pattern)
	}

	segments = strings.Split(pattern, "|")
	suffixes = strings.Split(segments[1], ",")

	suffixes = lo.Reject(suffixes, func(item string, _ int) bool {
		return item == ""
	})

	return segments, suffixes, nil
}

func newNodeFilter(def *core.FilterDef,
	fo *pref.FilterOptions,
) (core.TraverseFilter, error) {
	var (
		filter             core.TraverseFilter
		ifNotApplicable    = true
		err                error
		segments, suffixes []string
	)

	switch def.IfNotApplicable {
	case enums.TriStateBoolTrue:
		ifNotApplicable = true

	case enums.TriStateBoolFalse:
		ifNotApplicable = false

	case enums.TriStateBoolUndefined:
	}

	switch def.Type {
	case enums.FilterTypeExtendedGlob:
		if segments, suffixes, err = fromExtendedGlobPattern(def.Pattern); err != nil {
			return nil, err
		}

		base, exclusion := splitGlob(segments[0])

		filter = &ExtendedGlobFilter{
			Filter: Filter{
				name:            def.Description,
				scope:           def.Scope,
				pattern:         def.Pattern,
				negate:          def.Negate,
				ifNotApplicable: ifNotApplicable,
			},
			baseGlob: base,
			suffixes: lo.Map(suffixes, func(s string, _ int) string {
				return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
			}),
			anyExtension: slices.Contains(suffixes, "*"),
			exclusion:    exclusion,
		}

	case enums.FilterTypeRegex:
		filter = &RegexFilter{
			Filter: Filter{
				name:            def.Description,
				scope:           def.Scope,
				pattern:         def.Pattern,
				negate:          def.Negate,
				ifNotApplicable: ifNotApplicable,
			},
		}

	case enums.FilterTypeGlob:
		filter = &GlobFilter{
			Filter: Filter{
				name:            def.Description,
				scope:           def.Scope,
				pattern:         def.Pattern,
				negate:          def.Negate,
				ifNotApplicable: ifNotApplicable,
			},
		}

	case enums.FilterTypeCustom:
		if fo.Custom == nil {
			return nil, locale.ErrMissingCustomFilterDefinition
		}
		filter = fo.Custom

	case enums.FilterTypePoly:
		var polyE error

		if filter, polyE = newPolyFilter(fo.Node.Poly); polyE != nil {
			return nil, polyE
		}

	case enums.FilterTypeUndefined:
		return nil, locale.ErrFilterMissingType
	}

	if filter != nil {
		err = filter.Validate()
	}

	return filter, err
}

func newPolyFilter(polyDef *core.PolyFilterDef) (core.TraverseFilter, error) {
	// enforce the correct filter scopes
	//
	polyDef.File.Scope.Set(enums.ScopeFile)
	polyDef.File.Scope.Clear(enums.ScopeFolder)

	polyDef.Folder.Scope.Set(enums.ScopeFolder)
	polyDef.Folder.Scope.Clear(enums.ScopeFile)

	var (
		file, folder   core.TraverseFilter
		fileE, folderE error
	)

	if file, fileE = newNodeFilter(&polyDef.File, nil); fileE != nil {
		return nil, fileE
	}

	if folder, folderE = newNodeFilter(&polyDef.Folder, nil); folderE != nil {
		return nil, folderE
	}

	filter := &PolyFilter{
		File:   file,
		Folder: folder,
	}

	return filter, nil
}

const (
	exclusionDelim = "/"
)

func splitGlob(baseGlob string) (base, exclusion string) {
	base = strings.ToLower(baseGlob)

	if strings.Contains(base, exclusionDelim) {
		constituents := strings.Split(base, exclusionDelim)
		base = constituents[0]
		exclusion = constituents[1]
	}

	return base, exclusion
}

func newChildFilter(def *core.ChildFilterDef) (core.ChildTraverseFilter, error) {
	var (
		filter core.ChildTraverseFilter
	)

	if def == nil {
		return nil, locale.ErrFilterIsNil
	}

	switch def.Type {
	case enums.FilterTypeExtendedGlob:
		var (
			err                error
			segments, suffixes []string
		)

		if segments, suffixes, err = fromExtendedGlobPattern(def.Pattern); err != nil {
			return nil, locale.NewInvalidIncaseFilterDefError(def.Pattern)
		}

		base, exclusion := splitGlob(segments[0])

		filter = &ChildExtendedGlobFilter{
			ChildFilter: ChildFilter{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
			baseGlob: base,
			suffixes: lo.Map(suffixes, func(s string, _ int) string {
				return strings.ToLower(strings.TrimPrefix(strings.TrimSpace(s), "."))
			}),
			anyExtension: slices.Contains(suffixes, "*"),
			exclusion:    exclusion,
		}

	case enums.FilterTypeRegex:
		filter = &ChildRegexFilter{
			ChildFilter: ChildFilter{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
		}

	case enums.FilterTypeGlob:
		filter = &ChildGlobFilter{
			ChildFilter: ChildFilter{
				Name:    def.Description,
				Pattern: def.Pattern,
				Negate:  def.Negate,
			},
		}

	case enums.FilterTypeCustom:
		return nil, locale.ErrFilterCustomNotSupported

	case enums.FilterTypeUndefined:
		return nil, locale.ErrFilterUndefined

	case enums.FilterTypePoly:
	}

	if filter != nil {
		if err := filter.Validate(); err != nil {
			return nil, err
		}
	}

	return filter, nil
}

func newSampleFilter(def *core.SampleFilterDef,
	so *pref.SamplingOptions,
) (core.SampleTraverseFilter, error) {
	var (
		filter core.SampleTraverseFilter
	)

	if def == nil {
		return nil, locale.ErrFilterIsNil
	}

	base := SampleFilter{
		Filter: Filter{
			name:    def.Description,
			scope:   def.Scope.Scrub(),
			pattern: def.Pattern,
			negate:  def.Negate,
		},
	}

	if base.scope.IsFile() && so.NoOf.Files == 0 {
		return nil, locale.ErrInvalidFileSamplingSpecMissingFiles
	}

	if base.scope.IsFolder() && so.NoOf.Folders == 0 {
		return nil, locale.ErrInvalidFolderSamplingSpecMissingFolders
	}

	switch def.Type {
	case enums.FilterTypeExtendedGlob:
	case enums.FilterTypeRegex:
		filter = &SampleRegexFilter{
			SampleFilter: base,
		}
	case enums.FilterTypeGlob:
		filter = &SampleGlobFilter{
			SampleFilter: base,
		}

	case enums.FilterTypeCustom:
		if def.Custom == nil {
			return nil, locale.ErrFilterIsNil
		}
		filter = def.Custom
	case enums.FilterTypePoly:
	case enums.FilterTypeUndefined:
		return nil, locale.ErrFilterMissingType
	}

	if filter != nil {
		if err := filter.Validate(); err != nil {
			return nil, err
		}
	}

	return filter, nil
}
