package persist_test

import (
	"os"
	"testing/fstest"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/li18ngo"
	lab "github.com/snivilised/traverse/internal/laboratory"
	"github.com/snivilised/traverse/internal/opts"
	"github.com/snivilised/traverse/internal/persist"
	"github.com/snivilised/traverse/lfs"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Convert Options via JSON", Ordered, func() {
	var (
		FS lfs.TraverseFS
	)

	BeforeAll(func() {
		Expect(li18ngo.Use()).To(Succeed())
	})

	BeforeEach(func() {
		FS = &lab.TestTraverseFS{
			MapFS: fstest.MapFS{
				home: &fstest.MapFile{
					Mode: os.ModeDir,
				},
			},
		}

		_ = FS.MakeDirAll(destination, lab.Perms.Dir|os.ModeDir)
	})

	Context("ToJSON", func() {
		Context("given: source Options instance", func() {
			It("should: convert to JSON", func() {
				o, _, err := opts.Get(
					pref.WithDepth(4),
				)
				Expect(err).To(Succeed())
				Expect(persist.ToJSON(o)).To(HaveMarshaledEqual(o))
			})
		})
	})
})
