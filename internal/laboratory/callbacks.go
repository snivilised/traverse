package lab

import (
	"fmt"
	"strings"

	. "github.com/onsi/ginkgo/v2" //nolint:revive,stylecheck // ok
	. "github.com/onsi/gomega"    //nolint:revive,stylecheck // ok
	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/internal/third/lo"
)

func Begin(em string) cycle.BeginHandler {
	return func(state *cycle.BeginState) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:BEGIN], tree: '%v'\n", em, state.Tree,
		)
	}
}

func End(em string) cycle.EndHandler {
	return func(result core.TraverseResult) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:END], err: '%v'\n", em, result.Error(),
		)
	}
}

func UniversalCallback(name string) core.Client {
	return func(servant tv.Servant) error {
		node := servant.Node()
		depth := node.Extension.Depth
		GinkgoWriter.Printf(
			"---> 🌊 UNIVERSAL//%v-CALLBACK: (depth:%v) '%v'\n", name, depth, node.Path,
		)
		Expect(node.Extension).NotTo(BeNil(), Reason(node.Path))

		return nil
	}
}

func FoldersCallback(name string) core.Client {
	return func(servant tv.Servant) error {
		node := servant.Node()
		depth := node.Extension.Depth
		actualNoChildren := len(node.Children)
		GinkgoWriter.Printf(
			"---> 🔆 FOLDERS//CALLBACK%v: (depth:%v, children:%v) '%v'\n",
			name, depth, actualNoChildren, node.Path,
		)
		Expect(node.IsDirectory()).To(BeTrue(),
			Because(node.Path, "node expected to be folder"),
		)
		Expect(node.Extension).NotTo(BeNil(), Reason(node.Path))

		return nil
	}
}

func FilesCallback(name string) core.Client {
	return func(servant tv.Servant) error {
		node := servant.Node()
		GinkgoWriter.Printf("---> 🌙 FILES//%v-CALLBACK: '%v'\n", name, node.Path)
		Expect(node.IsDirectory()).To(BeFalse(),
			Because(node.Path, "node expected to be file"),
		)
		Expect(node.Extension).NotTo(BeNil(), Reason(node.Path))

		return nil
	}
}

func FoldersCaseSensitiveCallback(first, second string) core.Client {
	recording := make(RecordingMap)

	return func(servant tv.Servant) error {
		node := servant.Node()
		recording[node.Path] = len(node.Children)

		GinkgoWriter.Printf("---> 🔆 CASE-SENSITIVE-CALLBACK: '%v'\n", node.Path)
		Expect(node.IsDirectory()).To(BeTrue())

		if strings.HasSuffix(node.Path, second) {
			GinkgoWriter.Printf("---> 💧 FIRST: '%v', 💧 SECOND: '%v'\n", first, second)

			paths := lo.Keys(recording)
			_, found := lo.Find(paths, func(s string) bool {
				return strings.HasSuffix(s, first)
			})

			Expect(found).To(BeTrue(), fmt.Sprintf("for node: '%v'", node.Extension.Name))
		}

		return nil
	}
}
