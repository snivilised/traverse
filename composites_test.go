package age_test

import (
	"context"
	"fmt"
	"sync"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	age "github.com/snivilised/agenor"
	lab "github.com/snivilised/agenor/internal/laboratory"
	"github.com/snivilised/agenor/internal/services"
	"github.com/snivilised/agenor/locale"
	"github.com/snivilised/agenor/pref"
	"github.com/snivilised/li18ngo"
)

var (
	primeFacade = &pref.Using{
		Subscription: age.SubscribeFiles,
		Head: pref.Head{
			Handler: noOpHandler,
		},
		Tree: "tree",
	}

	resumeFacade = &pref.Relic{
		Head: pref.Head{
			Handler: noOpHandler,
		},
		From:     "path-to-json-file",
		Strategy: age.ResumeStrategyFastward,
	}
)

var _ = Describe("Composites", Ordered, func() {
	var (
		wg sync.WaitGroup
	)

	BeforeAll(func() {
		Expect(li18ngo.Use(
			func(o *li18ngo.UseOptions) {
				o.From.Sources = li18ngo.TranslationFiles{
					locale.SourceID: li18ngo.TranslationSource{Name: "agenor"},
				}
			},
		)).To(Succeed())
	})

	BeforeEach(func() {
		services.Reset()
	})

	// The point of these tests is not to check the validity of the navigation,
	// rather the point is just to make sure that the Composites can be invoked.
	// As long as there are no panics, we're happy; this is why there are no
	// expectations and we ignore the result and error. This enables us not to have
	// to arrange valid navigation state, which inevitably means that the actual
	// result will in many cases not be valid and the error therefore also may be
	// none nil.

	DescribeTable("hydra",
		func(specCtx SpecContext, entry *lab.CompositeTE) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, _ = age.Hydra(
				func() bool {
					return entry.IsWalk
				},
				func() bool {
					return entry.IsPrime
				},
				&wg,
			)(entry.Facade).Navigate(ctx)
		},
		func(entry *lab.CompositeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.Given, entry.Should,
			)
		},
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is walk/prime",
				Should: "return prime extent with sequential sync",
			},
			IsWalk:  true,
			IsPrime: true,
			Facade:  primeFacade,
		}),
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is walk/resume",
				Should: "return resume extent with sequential sync",
			},
			IsWalk:  true,
			IsPrime: false,
			Facade:  resumeFacade,
		}),
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is run/prime",
				Should: "return prime extent with concurrent sync",
			},
			IsWalk:  false,
			IsPrime: true,
			Facade:  primeFacade,
		}),
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is run/resume",
				Should: "return resume extent with concurrent sync",
			},
			IsWalk:  false,
			IsPrime: false,
			Facade:  resumeFacade,
		}),
	)

	DescribeTable("hare",
		func(specCtx SpecContext, entry *lab.CompositeTE) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, _ = age.Hare(func() bool {
				return entry.IsPrime
			}, &wg)(entry.Facade).Navigate(ctx)
		},
		func(entry *lab.CompositeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.Given, entry.Should,
			)
		},
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is prime",
				Should: "return prime extent with concurrent sync",
			},
			IsPrime: true,
			Facade:  primeFacade,
		}),
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is resume",
				Should: "return resume extent with concurrent sync",
			},
			IsPrime: false,
			Facade:  resumeFacade,
		}),
	)

	DescribeTable("tortoise",
		func(specCtx SpecContext, entry *lab.CompositeTE) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, _ = age.Tortoise(func() bool {
				return entry.IsPrime
			})(entry.Facade).Navigate(ctx)
		},
		func(entry *lab.CompositeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.Given, entry.Should,
			)
		},
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is prime",
				Should: "return prime extent with sequential sync",
			},
			IsPrime: true,
			Facade:  primeFacade,
		}),
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is resume",
				Should: "return resume extent with sequential sync",
			},
			IsPrime: false,
			Facade:  resumeFacade,
		}),
	)

	DescribeTable("goldfish",
		func(specCtx SpecContext, entry *lab.CompositeTE) {
			defer leaktest.Check(GinkgoT())()

			ctx, cancel := context.WithCancel(specCtx)
			defer cancel()

			_, _ = age.Goldfish(func() bool {
				return entry.IsWalk
			}, &wg)(entry.Facade).Navigate(ctx)
		},
		func(entry *lab.CompositeTE) string {
			return fmt.Sprintf("🧪 ===> given: '%v', should: '%v'",
				entry.Given, entry.Should,
			)
		},
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is walk",
				Should: "return prime extent with sequential sync",
			},
			IsWalk: true,
			Facade: primeFacade,
		}),
		Entry(nil, &lab.CompositeTE{
			AsyncTE: lab.AsyncTE{
				Given:  "is run",
				Should: "return prime extent with concurrent sync",
			},
			IsWalk: false,
			Facade: primeFacade,
		}),
	)
})
