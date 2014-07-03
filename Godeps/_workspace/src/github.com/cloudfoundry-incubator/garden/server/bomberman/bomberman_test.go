package bomberman_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"time"

	"github.com/cloudfoundry-incubator/garden/server/bomberman"
	"github.com/cloudfoundry-incubator/garden/warden"
	"github.com/cloudfoundry-incubator/garden/warden/fake_backend"
)

var _ = Describe("Bomberman", func() {
	It("straps a bomb to the given container with the container's grace time as the countdown", func() {
		detonated := make(chan warden.Container)

		bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
			detonated <- container
		})

		container := fake_backend.NewFakeContainer(
			warden.ContainerSpec{
				GraceTime: 100 * time.Millisecond,
				Handle:    "doomed",
			},
		)

		bomberman.Strap(container)

		select {
		case <-detonated:
		case <-time.After(container.GraceTime() + 50*time.Millisecond):
			Fail("did not detonate!")
		}
	})

	Context("when the container has a grace time of 0", func() {
		It("never detonates", func() {
			detonated := make(chan warden.Container)

			bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
				detonated <- container
			})

			container := fake_backend.NewFakeContainer(
				warden.ContainerSpec{
					GraceTime: 0,
					Handle:    "doomed",
				},
			)

			bomberman.Strap(container)

			select {
			case <-detonated:
				Fail("detonated!")
			case <-time.After(container.GraceTime() + 50*time.Millisecond):
			}
		})
	})

	Describe("pausing a container's timebomb", func() {
		It("prevents it from detonating", func() {
			detonated := make(chan warden.Container)

			bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
				detonated <- container
			})

			container := fake_backend.NewFakeContainer(
				warden.ContainerSpec{
					GraceTime: 100 * time.Millisecond,
					Handle:    "doomed",
				},
			)

			bomberman.Strap(container)
			bomberman.Pause("doomed")

			select {
			case <-detonated:
				Fail("detonated!")
			case <-time.After(container.GraceTime() + 50*time.Millisecond):
			}
		})

		Context("when the handle is invalid", func() {
			It("doesn't launch any missiles or anything like that", func() {
				bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
					panic("dont call me")
				})

				bomberman.Pause("BOOM?!")
			})
		})

		Describe("and then unpausing it", func() {
			It("causes it to detonate after the countdown", func() {
				detonated := make(chan warden.Container)

				bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
					detonated <- container
				})

				container := fake_backend.NewFakeContainer(
					warden.ContainerSpec{
						GraceTime: 100 * time.Millisecond,
						Handle:    "doomed",
					},
				)

				bomberman.Strap(container)
				bomberman.Pause("doomed")

				before := time.Now()
				bomberman.Unpause("doomed")

				select {
				case <-detonated:
					Ω(time.Since(before)).Should(BeNumerically(">=", 100*time.Millisecond))
				case <-time.After(container.GraceTime() + 50*time.Millisecond):
					Fail("did not detonate!")
				}
			})

			Context("when the handle is invalid", func() {
				It("doesn't launch any missiles or anything like that", func() {
					bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
						panic("dont call me")
					})

					bomberman.Unpause("BOOM?!")
				})
			})
		})
	})

	Describe("defusing a container's timebomb", func() {
		It("prevents it from detonating", func() {
			detonated := make(chan warden.Container)

			bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
				detonated <- container
			})

			container := fake_backend.NewFakeContainer(
				warden.ContainerSpec{
					GraceTime: 100 * time.Millisecond,
					Handle:    "doomed",
				},
			)

			bomberman.Strap(container)
			bomberman.Defuse("doomed")

			select {
			case <-detonated:
				Fail("detonated!")
			case <-time.After(container.GraceTime() + 50*time.Millisecond):
			}
		})

		Context("when the handle is invalid", func() {
			It("doesn't launch any missiles or anything like that", func() {
				bomberman := bomberman.New(fake_backend.New(), func(container warden.Container) {
					panic("dont call me")
				})

				bomberman.Defuse("BOOM?!")
			})
		})
	})
})
