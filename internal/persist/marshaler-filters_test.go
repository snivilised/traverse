package persist_test

import (
	"fmt"
	"os"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/pref"
)

// 📚 NB: these create functions are required because it is vitally important
// that objects are created consistently so as to not accidentally break
// equality checks. They do this by enforcing a single source of truth.

// 🍑 NODE:
func createJSONFilterFromCore(def *core.FilterDef) *json.FilterDef {
	return &json.FilterDef{
		Type:            def.Type,
		Description:     def.Description,
		Pattern:         def.Pattern,
		Scope:           def.Scope,
		Negate:          def.Negate,
		IfNotApplicable: def.IfNotApplicable,
	}
}

func createJSONFilterFromCoreWithPoly(def *core.FilterDef,
	poly *json.PolyFilterDef,
) *json.FilterDef {
	result := createJSONFilterFromCore(def)
	result.Poly = poly

	return result
}

func createCoreFilterDefFromJSON(jdef *json.FilterDef) *core.FilterDef {
	return &core.FilterDef{
		Type:            jdef.Type,
		Description:     jdef.Description,
		Pattern:         jdef.Pattern,
		Scope:           jdef.Scope,
		Negate:          jdef.Negate,
		IfNotApplicable: jdef.IfNotApplicable,
	}
}

func createCoreFilterDefFromJSONWithPoly(def *json.FilterDef,
	poly *core.PolyFilterDef,
) *core.FilterDef {
	result := createCoreFilterDefFromJSON(def)
	result.Poly = poly

	return result
}

// 🍑 CHILD:
func createChildFromNode(def *core.FilterDef) *core.ChildFilterDef {
	return &core.ChildFilterDef{
		Type:        def.Type,
		Description: def.Description,
		Pattern:     def.Pattern,
		Negate:      def.Negate,
	}
}

func createJSONChildFilterFromCore(def *core.ChildFilterDef) *json.ChildFilterDef {
	return &json.ChildFilterDef{
		Type:        def.Type,
		Description: def.Description,
		Pattern:     def.Pattern,
		Negate:      def.Negate,
	}
}

// 🍑 SAMPLE:
func createSampleFromNode(def *core.FilterDef) *core.SampleFilterDef {
	return &core.SampleFilterDef{
		Type:        def.Type,
		Description: def.Description,
		Pattern:     def.Pattern,
		Scope:       def.Scope,
		Negate:      def.Negate,
	}
}

func createJSONSampleFilterFromCore(def *core.SampleFilterDef) *json.SampleFilterDef {
	return &json.SampleFilterDef{
		Type:        def.Type,
		Description: def.Description,
		Pattern:     def.Pattern,
		Scope:       def.Scope,
		Negate:      def.Negate,
	}
}

func createJSONSampleFilterFromCoreWithPoly(
	def *core.FilterDef,
) *core.SampleFilterDef {
	return &core.SampleFilterDef{
		Poly: &core.PolyFilterDef{
			File:   *def,
			Folder: *def,
		},
	}
}

func createJSONSampleFilterDefFromCore(def *core.SampleFilterDef) *json.SampleFilterDef {
	return &json.SampleFilterDef{
		Type:        def.Type,
		Description: def.Description,
		Pattern:     def.Pattern,
		Scope:       def.Scope,
		Negate:      def.Negate,
	}
}

func createJSONSampleFilterDefFromCoreWithPoly(def *core.SampleFilterDef,
	poly *json.PolyFilterDef,
) *json.SampleFilterDef {
	result := createJSONSampleFilterDefFromCore(def)
	result.Poly = poly

	return result
}

var _ = Describe("Marshaler", Ordered, func() {
	var (
		FS       lfs.TraverseFS
		readPath string

		// 🍑 NODE:
		//
		sourceNodeFilterDef   *core.FilterDef
		jsonNodeFilterDef     json.FilterDef
		jsonPolyNodeFilterDef json.FilterDef
		polyNodeFilterDef     *core.FilterDef

		// 🍑 CHILD:
		//
		sourceChildFilterDef *core.ChildFilterDef
		jsonChildFilterDef   *json.ChildFilterDef

		// 🍑 SAMPLE:
		//
		sourceSampleFilterDef   *core.SampleFilterDef
		sampleFilterDef         *core.SampleFilterDef
		jsonSampleFilterDef     *json.SampleFilterDef
		jsonSamplePolyFilterDef *json.SampleFilterDef
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())
		readPath = source + "/" + restoreFile
	})

	BeforeEach(func() {
		FS = &lab.TestTraverseFS{
			MapFS: fstest.MapFS{
				home: &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}

		Expect(FS.MakeDirAll(destination, lab.Perms.Dir|os.ModeDir)).To(Succeed())
		Expect(FS.MakeDirAll(source, lab.Perms.Dir|os.ModeDir)).To(Succeed())
		Expect(FS.WriteFile(readPath, content, lab.Perms.File)).To(Succeed())

		// 🍑 NODE:
		//
		sourceNodeFilterDef = &core.FilterDef{
			Type:            enums.FilterTypeGlob,
			Description:     "items without .flac suffix",
			Pattern:         flac,
			Scope:           enums.ScopeAll,
			Negate:          true,
			IfNotApplicable: enums.TriStateBoolTrue,
		}

		jsonNodeFilterDef = *createJSONFilterFromCore(sourceNodeFilterDef)
		jsonPolyNodeFilterDef = *createJSONFilterFromCoreWithPoly(
			sourceNodeFilterDef, &json.PolyFilterDef{
				File:   jsonNodeFilterDef,
				Folder: jsonNodeFilterDef,
			},
		)

		polyNodeFilterDef = createCoreFilterDefFromJSONWithPoly(
			&jsonNodeFilterDef, &core.PolyFilterDef{
				File:   *sourceNodeFilterDef,
				Folder: *sourceNodeFilterDef,
			},
		)

		// 🍑 CHILD:
		//
		sourceChildFilterDef = createChildFromNode(sourceNodeFilterDef)
		jsonChildFilterDef = createJSONChildFilterFromCore(sourceChildFilterDef)

		// 🍑 SAMPLE:
		//
		sourceSampleFilterDef = createSampleFromNode(sourceNodeFilterDef)
		sampleFilterDef = createJSONSampleFilterFromCoreWithPoly(sourceNodeFilterDef)
		jsonSampleFilterDef = createJSONSampleFilterFromCore(sourceSampleFilterDef)
		jsonSamplePolyFilterDef = &json.SampleFilterDef{
			Poly: &json.PolyFilterDef{
				File:   jsonNodeFilterDef,
				Folder: jsonNodeFilterDef,
			},
		}

	})

	Context("map-fs", func() {
		DescribeTable("marshal filter defs",
			func(entry *marshalTE) {
				// This looks a bit odd, but actually helps us to reduce
				// the amount of test code required.
				//
				// marshal tweaks the JSON state to enforce unequal error, but
				// the tweak invoked by marshal can be shared by unmarshal,
				// without having to invoke unmarshal specific functionality.
				// The result of marshal can be passed into unmarshal.
				//
				unmarshal(entry, FS, readPath, marshal(entry, FS))
			},
			func(entry *marshalTE) string {
				return fmt.Sprintf("given: %v, 🧪 should: marshal successfully", entry.given)
			},
			// 🍉 FilterOptions.Node
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:pref.Options",
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &json.FilterDef{
						Type:        enums.FilterTypeRegex,
						Description: foo,
						Pattern:     flac,
						Scope:       enums.ScopeFile,
					}
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: sourceNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Type",
				},
				checkerTE: &checkerTE{
					field:   "Type",
					checker: check[enums.FilterType],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: sourceNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonNodeFilterDef
					result.JO.Filter.Node.Type = enums.FilterTypeExtendedGlob
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Description",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: sourceNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonNodeFilterDef
					result.JO.Filter.Node.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Pattern",
				},
				checkerTE: &checkerTE{
					field:   "Pattern",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: sourceNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonNodeFilterDef
					result.JO.Filter.Node.Pattern = bar
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Scope",
				},
				checkerTE: &checkerTE{
					field:   "Scope",
					checker: check[enums.FilterScope],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: sourceNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonNodeFilterDef
					result.JO.Filter.Node.Scope = enums.ScopeFile
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Negate",
				},
				checkerTE: &checkerTE{
					field:   "Negate",
					checker: check[bool],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: sourceNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonNodeFilterDef
					result.JO.Filter.Node.Negate = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.IfNotApplicable",
				},
				checkerTE: &checkerTE{
					field:   "IfNotApplicable",
					checker: check[enums.TriStateBool],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: sourceNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonNodeFilterDef
					result.JO.Filter.Node.IfNotApplicable = enums.TriStateBoolFalse
				},
			}),

			// 🍉 FilterOptions.Node.Poly
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions.Node.Poly - nil:pref.Options",
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonPolyNodeFilterDef
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions.Node.Poly - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: polyNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Poly.File",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: polyNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonPolyNodeFilterDef
					result.JO.Filter.Node.Poly.File.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Node.Poly.Folder",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Node: polyNodeFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Node = &jsonPolyNodeFilterDef
					result.JO.Filter.Node.Poly.Folder.Description = foo
				},
			}),

			// 🍉 FilterOptions.Child
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:pref.Options",
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Child = &json.ChildFilterDef{
						Type:        enums.FilterTypeGlob,
						Description: "items without .flac suffix",
						Pattern:     flac,
						Negate:      true,
					}
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: sourceChildFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Child = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Type",
				},
				checkerTE: &checkerTE{
					field:   "Type",
					checker: check[enums.FilterType],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: sourceChildFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Child = jsonChildFilterDef
					result.JO.Filter.Child.Type = enums.FilterTypeExtendedGlob
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Description",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: sourceChildFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Child = jsonChildFilterDef
					result.JO.Filter.Child.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Pattern",
				},
				checkerTE: &checkerTE{
					field:   "Pattern",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: sourceChildFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Child = jsonChildFilterDef
					result.JO.Filter.Child.Pattern = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Child.Negate",
				},
				checkerTE: &checkerTE{
					field:   "Negate",
					checker: check[bool],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Child: sourceChildFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Child = jsonChildFilterDef
					result.JO.Filter.Child.Negate = false
				},
			}),

			// 🍉 FilterOptions.Sample
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:pref.Options",
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSampleFilterDef
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sourceSampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Type",
				},
				checkerTE: &checkerTE{
					field:   "Type",
					checker: check[enums.FilterType],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sourceSampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSampleFilterDef
					result.JO.Filter.Sample.Type = enums.FilterTypeExtendedGlob
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Description",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sourceSampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSampleFilterDef
					result.JO.Filter.Sample.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Pattern",
				},
				checkerTE: &checkerTE{
					field:   "Pattern",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sourceSampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSampleFilterDef
					result.JO.Filter.Sample.Pattern = bar
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Scope",
				},
				checkerTE: &checkerTE{
					field:   "Scope",
					checker: check[enums.FilterScope],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sourceSampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSampleFilterDef
					result.JO.Filter.Sample.Scope = enums.ScopeFile
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Negate",
				},
				checkerTE: &checkerTE{
					field:   "Negate",
					checker: check[bool],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sourceSampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSampleFilterDef
					result.JO.Filter.Sample.Negate = false
				},
			}),

			// 🍉 FilterOptions.Sample.Poly
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions.Sample.Poly - nil:pref.Options",
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSamplePolyFilterDef
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions.Sample.Poly - nil:json.Options",
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = nil
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Poly.File",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSamplePolyFilterDef
					result.JO.Filter.Sample.Poly.File.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "FilterOptions - Sample.Poly.Folder",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithFilter(&pref.FilterOptions{
						Sample: sampleFilterDef,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Filter.Sample = jsonSamplePolyFilterDef
					result.JO.Filter.Sample.Poly.Folder.Description = foo
				},
			}),
		)
	})
})
