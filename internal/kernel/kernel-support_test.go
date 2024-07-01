package kernel_test

import (
	"io/fs"

	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/cycle"
	"github.com/snivilised/traverse/enums"
	"github.com/snivilised/traverse/internal/helpers"
)

const (
	RootPath    = "traversal-root-path"
	RestorePath = "/from-restore-path"
)

type recordingMap map[string]int
type recordingScopeMap map[string]enums.FilterScope
type recordingOrderMap map[string]int

type directoryQuantities struct {
	files    uint
	folders  uint
	children map[string]int
}

type naviTE struct {
	message       string
	should        string
	relative      string
	once          bool
	visit         bool
	caseSensitive bool
	subscription  enums.Subscription
	callback      core.Client
	mandatory     []string
	prohibited    []string
	expectedNoOf  directoryQuantities
}

func begin(em string) cycle.BeginHandler {
	return func(root string) {
		GinkgoWriter.Printf(
			"---> %v [traverse-navigator-test:BEGIN], root: '%v'\n", em, root,
		)
	}
}

func universalCallback(name string) core.Client {
	return func(node *core.Node) error {
		depth := node.Extension.Depth
		GinkgoWriter.Printf(
			"---> 🌊 UNIVERSAL//%v-CALLBACK: (depth:%v) '%v'\n", name, depth, node.Path,
		)
		Expect(node.Extension).NotTo(BeNil(), helpers.Reason(node.Path))
		return nil
	}
}

func subscribes(subscription enums.Subscription, de fs.DirEntry) bool {
	isAnySubscription := (subscription == enums.SubscribeUniversal)

	files := (subscription == enums.SubscribeFiles) && (!de.IsDir())
	folders := ((subscription == enums.SubscribeFolders) ||
		subscription == enums.SubscribeFoldersWithFiles) && (de.IsDir())

	return isAnySubscription || files || folders
}
