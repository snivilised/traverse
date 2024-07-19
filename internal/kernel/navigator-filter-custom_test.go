package kernel_test

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
	"github.com/snivilised/traverse/internal/lo"
	"github.com/snivilised/traverse/internal/services"
)

var _ = Describe("NavigatorFilterCustom", Ordered, func() {
	var (
		vfs  fstest.MapFS
		root string
	)

	BeforeAll(func() {
		const (
			verbose = true
		)

		vfs, root = helpers.Musico(verbose,
			filepath.Join("MUSICO", "RETRO-WAVE"),
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		services.Reset()
	})

	DescribeTable("custom-filter (glob)",
		func(ctx SpecContext, entry *filterTE) {
			recording := make(recordingMap)
			customFilter := &customFilter{
				name:    entry.name,
				pattern: entry.pattern,
				scope:   entry.scope,
				negate:  entry.negate,
			}

			path := helpers.Path(root, entry.relative)
			callback := func(item *core.Node) error {
				indicator := lo.Ternary(item.IsFolder(), "📁", "💠")
				GinkgoWriter.Printf(
					"===> %v Glob Filter(%v) source: '%v', item-name: '%v', item-scope(fs): '%v(%v)'\n",
					indicator,
					customFilter.Description(),
					customFilter.Source(),
					item.Extension.Name,
					item.Extension.Scope,
					customFilter.Scope(),
				)
				if lo.Contains(entry.mandatory, item.Extension.Name) {
					Expect(item).Should(MatchCurrentCustomFilter(customFilter))
				}

				recording[item.Extension.Name] = len(item.Children)
				return nil
			}
			result, err := tv.Walk().Configure().Extent(tv.Prime(
				&tv.Using{
					Root:         path,
					Subscription: entry.subscription,
					Handler:      callback,
					GetReadDirFS: func() fs.ReadDirFS {
						return vfs
					},
					GetQueryStatusFS: func(_ fs.FS) fs.StatFS {
						return vfs
					},
				},
				tv.WithFilterCustom(customFilter),
				tv.WithHookQueryStatus(
					func(qsys fs.StatFS, path string) (fs.FileInfo, error) {
						return qsys.Stat(helpers.TrimRoot(path))
					},
				),
				tv.WithHookReadDirectory(
					func(rfs fs.ReadDirFS, dirname string) ([]fs.DirEntry, error) {
						return rfs.ReadDir(helpers.TrimRoot(dirname))
					},
				),
			)).Navigate(ctx)

			assertNavigation(&entry.naviTE, testOptions{
				vfs:       vfs,
				recording: recording,
				path:      path,
				result:    result,
				err:       err,
			})
		},
		func(entry *filterTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v'", entry.given)
		},

		// === universal =====================================================

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): custom filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   8,
					folders: 0,
				},
			},
			name:    "items with '.flac' suffix",
			pattern: "*.flac",
			scope:   enums.ScopeAll,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(any scope): custom filter (negate)",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   6,
					folders: 8,
				},
			},
			name:    "items without .flac suffix",
			pattern: "*.flac",
			scope:   enums.ScopeAll,
			negate:  true,
		}),

		Entry(nil, &filterTE{
			naviTE: naviTE{
				given:        "universal(undefined scope): custom filter",
				relative:     "RETRO-WAVE",
				subscription: enums.SubscribeUniversal,
				expectedNoOf: quantities{
					files:   8,
					folders: 0,
				},
			},
			name:    "items with '.flac' suffix",
			pattern: "*.flac",
		}),
	)
})