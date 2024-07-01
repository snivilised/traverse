package tv_test

import (
	"context"
	"sync"

	"github.com/fortytw2/leaktest"
	. "github.com/onsi/ginkgo/v2" //nolint:revive // ok
	. "github.com/onsi/gomega"    //nolint:revive // ok

	tv "github.com/snivilised/traverse"
	"github.com/snivilised/traverse/core"
	"github.com/snivilised/traverse/internal/services"
	"github.com/snivilised/traverse/pref"
)

var _ = Describe("Director(Resume)", Ordered, func() {
	var restore pref.Option

	BeforeAll(func() {
		restore = func(o *tv.Options) error {
			o.Events.Begin.On(func(_ string) {})

			return nil
		}
	})

	BeforeEach(func() {
		services.Reset()
	})

	Context("simple", func() {
		Context("Walk", func() {
			It("🧪 should: walk resume navigation successfully", func(specCtx SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				const depth = 2

				_, err := tv.Walk().Configure().Extent(tv.Resume(
					&tv.Was{
						Using: tv.Using{
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						From:     RestorePath,
						Strategy: tv.ResumeStrategyFastward,
					},
					tv.WithDepth(depth),
					tv.WithOnDescend(func(_ *core.Node) {}),
					restore,
				)).Navigate(ctx)

				Expect(err).To(Succeed())
			})
		})

		Context("Run", func() {
			It("🧪 should: perform run navigation successfully", func(specCtx SpecContext) {
				defer leaktest.Check(GinkgoT())()

				ctx, cancel := context.WithCancel(specCtx)
				defer cancel()

				var wg sync.WaitGroup

				_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
					&tv.Was{
						Using: tv.Using{
							Subscription: tv.SubscribeFiles,
							Handler:      noOpHandler,
						},
						From:     RestorePath,
						Strategy: tv.ResumeStrategySpawn,
					},
					tv.WithOnDescend(func(_ *core.Node) {}),
					restore,
				)).Navigate(ctx)

				wg.Wait()
				Expect(err).To(Succeed())
			})
		})
	})

	Context("features", func() {
		Context("Run", func() {
			When("filter", func() {
				FIt("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
						&tv.Was{
							Using: tv.Using{
								Subscription: tv.SubscribeFiles,
								Handler:      noOpHandler,
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithFilter(&core.FilterDef{}),
						restore,
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})

			When("hibernate", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
						&tv.Was{
							Using: tv.Using{
								Subscription: tv.SubscribeFiles,
								Handler:      noOpHandler,
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithHibernationWake(&core.FilterDef{}),
						restore,
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})

			When("sample", func() {
				It("🧪 should: register ok", func(specCtx SpecContext) {
					defer leaktest.Check(GinkgoT())()

					ctx, cancel := context.WithCancel(specCtx)
					defer cancel()

					var wg sync.WaitGroup

					_, err := tv.Run(&wg).Configure().Extent(tv.Resume(
						&tv.Was{
							Using: tv.Using{
								Subscription: tv.SubscribeFiles,
								Handler:      noOpHandler,
							},
							From:     RestorePath,
							Strategy: tv.ResumeStrategySpawn,
						},
						tv.WithSampling(files, folders),
						tv.WithSampler(&pref.SamplerOptions{
							Iteration: pref.SamplingIterationOptions{
								Each:  func(_ *core.Node) bool { return false },
								While: func(_ *pref.FilteredInfo) bool { return false },
							},
						}),
						restore,
					)).Navigate(ctx)

					wg.Wait()
					Expect(err).To(Succeed())
				})
			})
		})
	})
})
