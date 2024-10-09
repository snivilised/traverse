package lfs_test

import (
	"fmt"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
)

var _ = Describe("file systems", Ordered, func() {
	var root string

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())

		root = lab.Repo("test")
	})

	Context("fs: ExistsInFS", func() {
		var fS lfs.ExistsInFS

		BeforeEach(func() {
			fS = lfs.NewExistsInFS(root)
		})

		Context("op: FileExists", func() {
			When("given: existing path", func() {
				It("🧪 should: return true", func() {
					Expect(AsFile(lab.Static.FS.Existing.File)).To(ExistInFS(fS))
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: return false", func() {
					Expect(AsFile(lab.Static.Foo)).NotTo(ExistInFS(fS))
				})
			})
		})

		Context("op: DirectoryExists", func() {
			When("given: existing path", func() {
				It("🧪 should: return true", func() {
					Expect(AsDirectory(lab.Static.FS.Existing.Directory)).To(ExistInFS(fS))
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: return false", func() {
					Expect(AsDirectory(lab.Static.Foo)).NotTo(ExistInFS(fS))
				})
			})
		})
	})

	Context("fs: ReadFileFS", func() {
		var fS lfs.ReadFileFS

		BeforeEach(func() {
			fS = lfs.NewReadFileFS(root)
		})

		Context("op: ReadFile", func() {
			When("given: existing path", func() {
				It("🧪 should: ", func() {
					_, err := fS.ReadFile(lab.Static.FS.Existing.File)
					Expect(err).To(Succeed())
				})
			})

			When("given: path does not exist", func() {
				It("🧪 should: ", func() {
					_, err := fS.ReadFile(lab.Static.Foo)
					Expect(err).NotTo(Succeed())
				})
			})
		})
	})

	Context("fs: MakeDirFS", func() {
		var (
			fS lfs.MakeDirFS
		)

		BeforeEach(func() {
			fS = lfs.NewMakeDirFS(root, false)
			scratchPath := filepath.Join(root, lab.Static.FS.Scratch)

			if _, err := os.Stat(scratchPath); err == nil {
				Expect(os.RemoveAll(scratchPath)).To(Succeed(),
					fmt.Sprintf("failed to delete existing directory %q", scratchPath),
				)
			}
		})

		Context("op: MakeDir", func() {
			When("given: path does not exist", func() {
				It("🧪 should: complete ok", func() {
					path := lab.Static.FS.Scratch
					Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
						Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
					)

					Expect(AsDirectory(path)).To(ExistInFS(fS))
				})
			})

			When("given: path already exists", func() {
				It("🧪 should: complete ok", func() {
					path := lab.Static.FS.Existing.Directory
					Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
						Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
					)
				})
			})
		})

		Context("op: MakeDirAll", func() {
			When("given: path does not exist", func() {
				It("🧪 should: complete ok", func() {
					path := lab.Static.FS.MakeDir.MakeAll
					Expect(fS.MakeDirAll(path, lab.Perms.Dir.Perm())).To(
						Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
					)

					Expect(AsDirectory(path)).To(ExistInFS(fS))
				})
			})

			When("given: path already exists", func() {
				It("🧪 should: complete ok", func() {
					path := lab.Static.FS.Existing.Directory
					Expect(fS.MakeDir(path, lab.Perms.Dir.Perm())).To(
						Succeed(), fmt.Sprintf("failed to MakeDir %q", path),
					)

					Expect(AsDirectory(path)).To(ExistInFS(fS))
				})
			})
		})
	})

	Context("fs: CopyFS", func() {
		var (
			fS     lfs.UniversalFS
			single string
		)

		BeforeEach(func() {
			single = filepath.Join(root, lab.Static.FS.Scratch)
		})

		BeforeEach(func() {
			fS = lfs.NewUniversalFS(root, false)
		})

		Context("op: Copy", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {
					_ = fS
					_ = single
				})
			})
		})

		Context("op: CopyAll", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {

				})
			})
		})
	})

	Context("fs: MoveFS", func() { // => table, lots of cases!
		Context("op: Move", func() {
			When("given: file", func() {
				It("🧪 should: ", func() {

				})
			})
		})
	})

	Context("fs: RemoveFS", func() {
		// file(exits/not), folder(exits/not)
		Context("op: Remove", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {

				})
			})
		})

		Context("op: RemoveAll", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {

				})
			})
		})
	})

	Context("fs: RenameFS", func() {
		Context("op: Rename", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {

				})
			})
		})
	})

	Context("fs: WriteFileFS", func() {
		Context("op: Create", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {

				})
			})
		})

		Context("op: WriteFile", func() {
			When("given: ", func() {
				It("🧪 should: ", func() {

				})
			})
		})
	})
})
