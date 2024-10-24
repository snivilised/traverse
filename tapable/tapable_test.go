package tapable_test

import (
	"io/fs"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/pref"
)

const (
	tree      = "/traversal-tree-path"
	spoofed   = "spoofed"
	respoofed = "re-spoofed"
	verbose   = false
)

var (
	fakeSubPath = &core.SubPathInfo{
		Tree: tree,
		Node: &core.Node{
			Extension: core.Top("/tree", nil).Extension,
		},
	}
)

var _ = Describe("Tapable", Ordered, func() {
	var (
		invoked bool
		o       *pref.Options
		err     error
		FS      *lab.TestTraverseFS
		root    string
	)

	BeforeAll(func() {
		FS, root = lab.Musico(verbose,
			lab.Static.RetroWave,
		)
		Expect(root).NotTo(BeEmpty())
	})

	BeforeEach(func() {
		invoked = false
		o, _, err = opts.Get()
		Expect(err).To(Succeed())
	})

	Context("hooks", func() {
		Context("FileSubPath", func() {
			Context("Chain", func() {
				When("single", func() {
					It("🧪 should: invoke", func() {
						o.Hooks.FileSubPath.Chain(
							func(_ string, _ *core.SubPathInfo) string {
								return spoofed
							},
						)
						result := o.Hooks.FileSubPath.Invoke()(fakeSubPath)

						Expect(result).To(Equal(spoofed), "FileSubPath hook not invoked")
					})
				})

				When("multiple", func() {
					It("🧪 should: broadcast", func() {
						o.Hooks.FileSubPath.Chain(
							func(_ string, _ *core.SubPathInfo) string {
								return spoofed
							},
						)
						o.Hooks.FileSubPath.Chain(
							func(_ string, _ *core.SubPathInfo) string {
								return respoofed
							},
						)
						result := o.Hooks.FileSubPath.Invoke()(fakeSubPath)

						Expect(result).To(Equal(respoofed), "FileSubPath hook not broadcasted")
					})
				})
			})

			When("Tap", func() {
				It("🧪 should: invoke hook", func() {
					o.Hooks.FileSubPath.Tap(
						func(_ *core.SubPathInfo) string {
							invoked = true
							return ""
						},
					)
					o.Hooks.FileSubPath.Default()(fakeSubPath)
					o.Hooks.FileSubPath.Invoke()(fakeSubPath)

					Expect(invoked).To(BeTrue(), "FileSubPath hook not invoked")
				})
			})
		})

		Context("FolderSubPath ", func() {
			Context("Chain", func() {
				When("single", func() {
					It("🧪 should: invoke", func() {
						o.Hooks.FolderSubPath.Chain(
							func(_ string, _ *core.SubPathInfo) string {
								return spoofed
							},
						)
						result := o.Hooks.FolderSubPath.Invoke()(fakeSubPath)

						Expect(result).To(Equal(spoofed), "FolderSubPath hook not invoked")
					})
				})

				When("multiple", func() {
					It("🧪 should: broadcast", func() {
						o.Hooks.FolderSubPath.Chain(func(_ string, _ *core.SubPathInfo) string {
							return spoofed
						})
						o.Hooks.FolderSubPath.Chain(func(_ string, _ *core.SubPathInfo) string {
							return respoofed
						})
						result := o.Hooks.FolderSubPath.Invoke()(fakeSubPath)

						Expect(result).To(Equal(respoofed), "FolderSubPath hook not invoked")
					})
				})
			})

			When("Tap", func() {
				It("🧪 should: invoke hook", func() {
					o.Hooks.FolderSubPath.Tap(func(_ *core.SubPathInfo) string {
						invoked = true
						return ""
					})
					o.Hooks.FolderSubPath.Default()(fakeSubPath)
					o.Hooks.FolderSubPath.Invoke()(fakeSubPath)

					Expect(invoked).To(BeTrue(), "FolderSubPath hook not invoked")
				})
			})

		})

		Context("ReadDirectory", func() {
			Context("Chain", func() {
				When("single", func() {
					It("🧪 should: invoke", func() {
						path := lab.Static.RetroWave
						o.Hooks.ReadDirectory.Chain(
							func(result []fs.DirEntry, err error,
								_ fs.ReadDirFS, _ string,
							) ([]fs.DirEntry, error) {
								return result, err
							},
						)

						result, err := o.Hooks.ReadDirectory.Invoke()(FS, path)
						Expect(err).To(Succeed())
						Expect(result).To(
							lab.HaveDirectoryContents(
								[]string{"Chromatics", "College", "Electric Youth"},
							),
							"ReadDirectory hook not invoked",
						)
					})
				})

				When("multiple", func() {
					It("🧪 should: broadcast", func() {
						path := lab.Static.RetroWave
						o.Hooks.ReadDirectory.Chain(
							func(result []fs.DirEntry, err error,
								_ fs.ReadDirFS, _ string,
							) ([]fs.DirEntry, error) {
								return result, err
							},
						)
						o.Hooks.ReadDirectory.Chain(
							func(result []fs.DirEntry, err error,
								_ fs.ReadDirFS, _ string,
							) ([]fs.DirEntry, error) {
								return []fs.DirEntry{result[0]}, err
							},
						)

						result, e := o.Hooks.ReadDirectory.Invoke()(FS, path)
						Expect(e).To(Succeed())
						Expect(result).To(
							lab.HaveDirectoryContents(
								[]string{"Chromatics"},
							),
							"ReadDirectory hook not broadcasted",
						)
					})
				})
			})

			When("Tap", func() {
				It("🧪 should: invoke hook", func() {
					o.Hooks.ReadDirectory.Tap(
						func(_ fs.ReadDirFS, _ string) ([]fs.DirEntry, error) {
							invoked = true
							return []fs.DirEntry{}, nil
						},
					)

					sys := tv.NewReadDirFS(tv.Rel{
						Root: root,
					})
					_, _ = o.Hooks.ReadDirectory.Default()(sys, root)
					_, _ = o.Hooks.ReadDirectory.Invoke()(sys, root)

					Expect(invoked).To(BeTrue(), "ReadDirectory hook not invoked")
				})
			})
		})

		Context("QueryStatus", func() {
			Context("Chain", func() {
				When("single", func() {
					It("🧪 should: invoke", func() {
						path := lab.Static.RetroWave
						o.Hooks.QueryStatus.Chain(
							func(result fs.FileInfo, err error,
								_ fs.StatFS, _ string,
							) (fs.FileInfo, error) {
								invoked = true
								return result, err
							},
						)
						_, err := o.Hooks.QueryStatus.Invoke()(FS, path)

						Expect(err).To(Succeed())
						Expect(invoked).To(BeTrue(), "QueryStatus hook not invoked")
					})
				})

				When("multiple", func() {
					It("🧪 should: broadcast", func() {
						path := lab.Static.RetroWave
						o.Hooks.QueryStatus.Chain(
							func(result fs.FileInfo, err error,
								_ fs.StatFS, _ string,
							) (fs.FileInfo, error) {
								return result, err
							},
						)
						o.Hooks.QueryStatus.Chain(
							func(result fs.FileInfo, err error,
								_ fs.StatFS, _ string,
							) (fs.FileInfo, error) {
								invoked = true
								return result, err
							},
						)
						_, e := o.Hooks.QueryStatus.Invoke()(FS, path)

						Expect(e).To(Succeed())
						Expect(invoked).To(BeTrue(), "QueryStatus hook not broadcasted")
					})
				})
			})

			When("Tap", func() {
				It("🧪 should: invoke hook", func() {
					o.Hooks.QueryStatus.Tap(
						func(_ fs.StatFS, _ string) (fs.FileInfo, error) {
							invoked = true
							return nil, nil
						},
					)
					_, _ = o.Hooks.QueryStatus.Default()(FS, root)
					_, _ = o.Hooks.QueryStatus.Invoke()(FS, root)

					Expect(invoked).To(BeTrue(), "QueryStatus hook not invoked")
				})
			})
		})

		Context("Sort", func() {
			Context("Chain", func() {
				When("single", func() {
					It("🧪 should: invoke", func() {
						o.Hooks.Sort.Chain(
							func(_ []fs.DirEntry, _ ...any) {
								invoked = true
							},
						)
						o.Hooks.Sort.Invoke()([]fs.DirEntry{})

						Expect(invoked).To(BeTrue(), "Sort hook not invoked")
					})
				})

				When("multiple", func() {
					It("🧪 should: broadcast", func() {
						o.Hooks.Sort.Chain(
							func(_ []fs.DirEntry, _ ...any) {},
						)
						o.Hooks.Sort.Chain(
							func(_ []fs.DirEntry, _ ...any) {
								invoked = true
							},
						)
						o.Hooks.Sort.Invoke()([]fs.DirEntry{})

						Expect(invoked).To(BeTrue(), "Sort hook not broadcasted")
					})
				})
			})

			When("Tap", func() {
				It("🧪 should: invoke hook", func() {
					o.Hooks.Sort.Tap(func(_ []fs.DirEntry, _ ...any) {
						invoked = true
					})
					o.Hooks.Sort.Default()([]fs.DirEntry{})
					o.Hooks.Sort.Invoke()([]fs.DirEntry{})

					Expect(invoked).To(BeTrue(), "Sort hook not invoked")
				})
			})
		})
	})
})
