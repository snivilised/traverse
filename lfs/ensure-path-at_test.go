package lfs_test

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/locale"
)

var _ = Describe("EnsurePathAt", Ordered, func() {
	var (
		mocks *lfs.ResolveMocks
		mfs   *makeDirMapFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "traverse"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		mocks = &lfs.ResolveMocks{
			HomeFunc: func() (string, error) {
				return filepath.Join(string(filepath.Separator), "home", "prodigy"), nil
			},
			AbsFunc: func(_ string) (string, error) {
				return "", errors.New("not required for these tests")
			},
		}

		mfs = &makeDirMapFS{
			mapFS: fstest.MapFS{
				filepath.Join("home", "prodigy"): &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}
	})

	DescribeTable("with mapFS",
		func(entry *ensureTE) {
			home, _ := mocks.HomeFunc()
			location := lab.TrimRoot(filepath.Join(home, entry.relative))

			if entry.directory {
				location += string(filepath.Separator)
			}

			actual, err := lfs.EnsurePathAt(location, "default-test.log", lab.Perms.File, mfs)
			directory, _ := filepath.Split(actual)
			directory = filepath.Clean(directory)
			expected := lab.TrimRoot(lab.Path(home, entry.expected))

			Expect(err).Error().To(BeNil())
			Expect(actual).To(Equal(expected))
			Expect(AsDirectory(lab.TrimRoot(directory))).To(ExistInFS(mfs))
		},
		func(entry *ensureTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'", entry.given, entry.should)
		},

		Entry(nil, &ensureTE{
			given:    "path is file",
			should:   "create parent directory and return specified file path",
			relative: filepath.Join("logs", "test.log"),
			expected: "logs/test.log",
		}),

		Entry(nil, &ensureTE{
			given:     "path is directory",
			should:    "create parent directory and return default file path",
			relative:  "logs/",
			directory: true,
			expected:  "logs/default-test.log",
		}),
	)
})
