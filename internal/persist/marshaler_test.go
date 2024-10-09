package persist_test

import (
	"errors"
	"fmt"
	"os"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/opts/json"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/internal/third/lo"
	"github.com/snivilised/traverse/internal/types"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/locale"
	"github.com/snivilised/traverse/pref"
)

func check[T any](entry *checkerTE, err error) error {
	if err, ok := errors.Unwrap(err).(persist.UnequalValueError[T]); ok {
		return lo.Ternary(err.Field == entry.field,
			nil, fmt.Errorf("actual(%q) => expected: %v; value: '%v', other: '%v'",
				err.Field, entry.field, err.Value, err.Other,
			),
		)
	}

	return &wrongUnequalError{
		field: entry.field,
		err:   err,
	}
}

func marshal(entry *marshalTE, tfs lfs.TraverseFS) *tampered {
	// success:
	o, _, err := opts.Get(
		pref.IfOptionF(entry.option != nil, func() pref.Option {
			return entry.option()
		}),
	)
	Expect(err).To(Succeed(), "MARSHAL")

	writePath := destination + "/" + tempFile
	request := &persist.MarshalRequest{
		O: o,
		Active: &types.ActiveState{
			Root:        destination,
			Hibernation: enums.HibernationPending,
			CurrentPath: "/top/a/b/c",
			Depth:       3,
		},
		Path: writePath,
		Perm: lab.Perms.File,
		FS:   tfs,
	}
	result, err := persist.Marshal(request)

	Expect(err).To(Succeed(), "MARSHAL")
	Expect(result).NotTo(BeNil())

	// unequal error:
	if entry.tweak != nil {
		entry.tweak(result)
	}

	e := persist.Equals(&persist.Comparison{
		O:  o,
		JO: result.JO,
	})

	Expect(e).NotTo(Succeed(), "MARSHAL")
	if e != nil && entry.checkerTE != nil && entry.checkerTE.checker != nil {
		Expect(entry.checker(entry.checkerTE, e)).To(Succeed(), "MARSHAL")
	}

	return &tampered{
		o:      o,
		result: result,
	}
}

func unmarshal(entry *marshalTE, tfs lfs.TraverseFS, restorePath string, t *tampered) {
	// success:
	request := &persist.UnmarshalRequest{
		Restore: &types.RestoreState{
			Path:   restorePath,
			FS:     tfs,
			Resume: enums.ResumeStrategySpawn,
		},
	}
	state, err := persist.Unmarshal(request, entry.tweak)
	Expect(err).To(Succeed(), "UNMARSHAL")

	// unequal error:
	e := persist.Equals(&persist.Comparison{
		O:  t.o,
		JO: state.JO,
	})

	Expect(e).NotTo(Succeed(), "UNMARSHAL")

	if e != nil && entry.checkerTE != nil && entry.checkerTE.checker != nil {
		Expect(entry.checker(entry.checkerTE, e)).To(Succeed(), "UNMARSHAL")
	}
}

func createJSONSamplingOptions(so *pref.SamplingOptions) *json.SamplingOptions {
	return &json.SamplingOptions{
		Type:      so.Type,
		InReverse: so.InReverse,
		NoOf: json.EntryQuantities{
			Files:   so.NoOf.Files,
			Folders: so.NoOf.Folders,
		},
	}
}

var _ = Describe("Marshaler", Ordered, func() {
	var (
		FS lfs.TraverseFS

		sourceNodeFilterDef *core.FilterDef
		jsonNodeFilterDef   json.FilterDef
		samplingOptions     *pref.SamplingOptions
		jsonSamplingOptions *json.SamplingOptions

		readPath string
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

		sourceNodeFilterDef = &core.FilterDef{
			Type:            enums.FilterTypeGlob,
			Description:     "items without .flac suffix",
			Pattern:         flac,
			Scope:           enums.ScopeAll,
			Negate:          true,
			IfNotApplicable: enums.TriStateBoolTrue,
		}

		jsonNodeFilterDef = *createJSONFilterFromCore(sourceNodeFilterDef)

		samplingOptions = &pref.SamplingOptions{
			Type:      enums.SampleTypeFilter,
			InReverse: true,
			NoOf: pref.EntryQuantities{
				Files:   2,
				Folders: 3,
			},
		}

		jsonSamplingOptions = createJSONSamplingOptions(samplingOptions)
	})

	Context("map-fs", func() {
		DescribeTable("marshal",
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

			// 🍉 NavigationBehaviours:
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SubPathBehaviour.KeepTrailingSep",
				},
				checkerTE: &checkerTE{
					field:   "KeepTrailingSep",
					checker: check[bool],
				},
				option: func() pref.Option {
					return pref.WithSubPathBehaviour(&pref.SubPathBehaviour{
						KeepTrailingSep: false,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Behaviours.SubPath.KeepTrailingSep = true
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.WithSortBehaviour.IsCaseSensitive",
				},
				checkerTE: &checkerTE{
					field:   "IsCaseSensitive",
					checker: check[bool],
				},
				option: func() pref.Option {
					return pref.WithSortBehaviour(&pref.SortBehaviour{
						IsCaseSensitive: true,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Behaviours.Sort.IsCaseSensitive = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.WithSortBehaviour.SortFilesFirst",
				},
				checkerTE: &checkerTE{
					field:   "SortFilesFirst",
					checker: check[bool],
				},
				option: func() pref.Option {
					return pref.WithSortBehaviour(&pref.SortBehaviour{
						SortFilesFirst: true,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Behaviours.Sort.SortFilesFirst = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.CascadeBehaviour.WithDepth",
				},
				checkerTE: &checkerTE{
					field:   "Depth",
					checker: check[uint],
				},
				option: func() pref.Option {
					return pref.WithDepth(4)
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Behaviours.Cascade.Depth = 99
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.CascadeBehaviour.NoRecurse",
				},
				checkerTE: &checkerTE{
					field:   "NoRecurse",
					checker: check[bool],
				},
				option: pref.WithNoRecurse,
				tweak: func(result *persist.MarshalResult) {
					result.JO.Behaviours.Cascade.NoRecurse = false
				},
			}),

			// 🍉 SamplingOptions:
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.Type",
				},
				checkerTE: &checkerTE{
					field:   "Type",
					checker: check[enums.SampleType],
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(samplingOptions)
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Sampling = *jsonSamplingOptions
					result.JO.Sampling.Type = enums.SampleTypeSlice
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.InReverse",
				},
				checkerTE: &checkerTE{
					field:   "InReverse",
					checker: check[bool],
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(&pref.SamplingOptions{
						InReverse: true,
					})
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Sampling.InReverse = false
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.NoOf.Files",
				},
				checkerTE: &checkerTE{
					field:   "Files",
					checker: check[uint],
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(samplingOptions)
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Sampling = *jsonSamplingOptions
					result.JO.Sampling.NoOf.Files = 99
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "NavigationBehaviours.SamplingOptions.NoOf.Folders",
				},
				option: func() pref.Option {
					return pref.WithSamplingOptions(samplingOptions)
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Sampling.NoOf.Folders = 99
				},
			}),

			// 🍉 HibernateOptions:
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveWake",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithHibernationFilterWake(sourceNodeFilterDef)
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Hibernate.WakeAt = &jsonNodeFilterDef
					result.JO.Hibernate.WakeAt.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveSleep",
				},
				checkerTE: &checkerTE{
					field:   "Description",
					checker: check[string],
				},
				option: func() pref.Option {
					return pref.WithHibernationFilterSleep(sourceNodeFilterDef)
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Hibernate.SleepAt = &jsonNodeFilterDef
					result.JO.Hibernate.SleepAt.Description = foo
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveWake",
				},
				checkerTE: &checkerTE{
					field:   "InclusiveWake",
					checker: check[bool],
				},
				option: pref.WithHibernationBehaviourExclusiveWake,
				tweak: func(result *persist.MarshalResult) {
					result.JO.Hibernate.Behaviour.InclusiveWake = true
				},
			}),

			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "HibernateOptions.Behaviour.InclusiveSleep",
				},
				checkerTE: &checkerTE{
					field:   "InclusiveSleep",
					checker: check[bool],
				},
				option: pref.WithHibernationBehaviourInclusiveSleep,
				tweak: func(result *persist.MarshalResult) {
					result.JO.Hibernate.Behaviour.InclusiveSleep = false
				},
			}),

			// 🍉 ConcurrencyOptions:
			//
			Entry(nil, &marshalTE{
				persistTE: persistTE{
					given: "ConcurrencyOptions.NoW",
				},
				checkerTE: &checkerTE{
					field:   "NoW",
					checker: check[uint],
				},
				option: func() pref.Option {
					return pref.WithNoW(5)
				},
				tweak: func(result *persist.MarshalResult) {
					result.JO.Concurrency.NoW = 99
				},
			}),
		)

		Context("UnequalPtrError", func() {
			When("pref.Options is nil", func() {
				It("🧪 should: return UnequalPtrError", func() {
					err := persist.Equals(&persist.Comparison{
						JO: &json.Options{},
					})

					Expect(err).NotTo(Succeed())
					Expect(errors.Is(err, locale.ErrUnEqualConversion)).To(BeTrue(),
						"error should be a locale.ErrUnEqualConversion",
					)
				})
			})

			When("json FilterDef is nil", func() {
				It("🧪 should: return UnequalPtrError", func() {
					o, _, _ := opts.Get()

					err := persist.Equals(&persist.Comparison{
						O: o,
					})

					Expect(err).NotTo(Succeed())
					Expect(errors.Is(err, locale.ErrUnEqualConversion)).To(BeTrue(),
						"error should be a locale.ErrUnEqualConversion",
					)
				})
			})
		})
	})
})
