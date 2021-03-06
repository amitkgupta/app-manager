package handler_test

import (
	"errors"
	"syscall"

	. "github.com/cloudfoundry-incubator/app-manager/handler"
	"github.com/cloudfoundry-incubator/app-manager/start_message_builder"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/fake_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/cloudfoundry/storeadapter"
	"github.com/pivotal-golang/lager/lagertest"
	"github.com/tedsuo/ifrit"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Handler", func() {
	var (
		startMessageBuilder       *start_message_builder.StartMessageBuilder
		bbs                       *fake_bbs.FakeAppManagerBBS
		logger                    *lagertest.TestLogger
		desiredLRP                models.DesiredLRP
		repAddrRelativeToExecutor string
		healthChecks              map[string]string

		handler ifrit.Process
	)

	BeforeEach(func() {
		bbs = fake_bbs.NewFakeAppManagerBBS()

		repAddrRelativeToExecutor = "127.0.0.1:20515"

		healthChecks = map[string]string{
			"some-stack": "some-health-check.tgz",
		}

		logger = lagertest.NewTestLogger("test")

		startMessageBuilder = start_message_builder.New(repAddrRelativeToExecutor, healthChecks, logger)

		handlerRunner := NewHandler(bbs, startMessageBuilder, logger)

		desiredLRP = models.DesiredLRP{
			ProcessGuid:  "the-app-guid-the-app-version",
			Source:       "http://the-droplet.uri.com",
			Stack:        "some-stack",
			StartCommand: "the-start-command",
			Environment: []models.EnvironmentVariable{
				{Name: "foo", Value: "bar"},
				{Name: "VCAP_APPLICATION", Value: "{\"application_name\":\"my-app\"}"},
			},
			MemoryMB:        128,
			DiskMB:          512,
			FileDescriptors: 32,
			Instances:       2,
			Routes:          []string{"route1", "route2"},
			LogGuid:         "the-log-id",
		}
		handler = ifrit.Envoke(handlerRunner)
	})

	AfterEach(func(done Done) {
		handler.Signal(syscall.SIGINT)
		<-handler.Wait()
		Eventually(bbs.DesiredLRPStopChan).Should(BeClosed())
		close(done)
	})

	Describe("lifecycle", func() {
		Describe("waiting until all desired are processed before shutting down", func() {
			var receivedAuctions chan models.LRPStartAuction

			BeforeEach(func() {
				receivedAuctions = make(chan models.LRPStartAuction)
				bbs.WhenRequestingLRPStartAuctions = func(lrp models.LRPStartAuction) error {
					receivedAuctions <- lrp
					return nil
				}
			})

			It("should not shut down until all desireds are processed", func() {
				bbs.DesiredLRPChangeChan <- models.DesiredLRPChange{
					Before: nil,
					After:  &desiredLRP,
				}
				bbs.DesiredLRPChangeChan <- models.DesiredLRPChange{
					Before: nil,
					After:  &desiredLRP,
				}

				handler.Signal(syscall.SIGINT)
				didShutDown := handler.Wait()

				Consistently(didShutDown).ShouldNot(Receive())

				for i := 0; i < desiredLRP.Instances*2; i++ {
					Eventually(receivedAuctions).Should(Receive())
				}

				Eventually(didShutDown).Should(Receive())
			})
		})

		Describe("when an error occurs", func() {
			var newChan chan models.DesiredLRPChange
			BeforeEach(func() {
				newChan = make(chan models.DesiredLRPChange, 1)
				bbs.DesiredLRPChangeChan = newChan
				bbs.DesiredLRPErrChan <- errors.New("oops")
			})

			It("should reestablish the watch", func() {
				newChan <- models.DesiredLRPChange{
					Before: nil,
					After:  &desiredLRP,
				}

				Eventually(bbs.GetLRPStartAuctions).Should(HaveLen(2))
			})
		})

		Describe("when the desired channel is closed", func() {
			var newChan chan models.DesiredLRPChange
			BeforeEach(func() {
				newChan = make(chan models.DesiredLRPChange, 1)
				oldChan := bbs.DesiredLRPChangeChan
				bbs.DesiredLRPChangeChan = newChan
				close(oldChan)
			})

			It("should reestablish the watch", func() {
				newChan <- models.DesiredLRPChange{
					Before: nil,
					After:  &desiredLRP,
				}

				Eventually(bbs.GetLRPStartAuctions).Should(HaveLen(2))
			})
		})

	})

	Describe("when a desired LRP change message is received", func() {
		JustBeforeEach(func() {
			bbs.DesiredLRPChangeChan <- models.DesiredLRPChange{
				Before: nil,
				After:  &desiredLRP,
			}
		})

		Describe("the happy path", func() {
			BeforeEach(func() {
				bbs.WhenGettingAvailableFileServer = func() (string, error) {
					return "http://file-server.com/", nil
				}
			})

			It("puts a LRPStartAuction in the bbs", func() {
				Eventually(bbs.GetLRPStartAuctions).Should(HaveLen(2))

				startAuctions := bbs.GetLRPStartAuctions()

				firstStartAuction := startAuctions[0]
				Ω(firstStartAuction.ProcessGuid).Should(Equal("the-app-guid-the-app-version"))

				secondStartAuction := startAuctions[1]
				Ω(secondStartAuction.ProcessGuid).Should(Equal("the-app-guid-the-app-version"))
			})

			It("assigns increasing indices for the auction requests", func() {
				Eventually(bbs.GetLRPStartAuctions).Should(HaveLen(2))
				startAuctions := bbs.GetLRPStartAuctions()

				firstStartAuction := startAuctions[0]
				secondStartAuction := startAuctions[1]

				Ω(firstStartAuction.Index).Should(Equal(0))
				Ω(*firstStartAuction.Log.Index).Should(Equal(0))
				Ω(secondStartAuction.Index).Should(Equal(1))
				Ω(*secondStartAuction.Log.Index).Should(Equal(1))
			})
		})

		Context("when file server is not available", func() {
			BeforeEach(func() {
				bbs.WhenGettingAvailableFileServer = func() (string, error) {
					return "", storeadapter.ErrorKeyNotFound
				}
			})

			It("does not put a LRPStartAuction in the bbs", func() {
				Consistently(bbs.GetLRPStartAuctions).Should(BeEmpty())
			})
		})

		Context("when unable to build a start message", func() {
			BeforeEach(func() {
				desiredLRP.Stack = "some-unknown-stack"
			})

			It("does not put a LRPStartAuction in the bbs", func() {
				Consistently(bbs.GetLRPStartAuctions).Should(BeEmpty())
			})
		})

		Context("when there is an error writing a LRPStartAuction to the BBS", func() {
			BeforeEach(func() {
				bbs.LRPStartAuctionErr = errors.New("connection error")
			})

			It("logs an error", func() {
				Eventually(logger.TestSink.Buffer).Should(gbytes.Say("handler.desired-lrp-change.request-start-auction-failed"))
			})
		})

		Context("when there is an error fetching the actual instances", func() {
			BeforeEach(func() {
				bbs.ActualLRPsErr = errors.New("connection error")
			})

			It("does not put a LRPStartAuction in the bbs", func() {
				Consistently(bbs.GetLRPStartAuctions).Should(BeEmpty())
			})
		})

		Context("when there are already instances running for the desired app, but some are missing", func() {
			BeforeEach(func() {
				desiredLRP.Instances = 4
				bbs.Lock()
				bbs.ActualLRPs = []models.ActualLRP{
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "a",
						Index:        0,
						State:        models.ActualLRPStateStarting,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "b",
						Index:        4,
						State:        models.ActualLRPStateRunning,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "c",
						Index:        5,
						State:        models.ActualLRPStateRunning,
					},
				}
				bbs.Unlock()
			})

			It("only starts missing ones", func() {
				Eventually(bbs.GetLRPStartAuctions).Should(HaveLen(3))
				startAuctions := bbs.GetLRPStartAuctions()

				Ω(startAuctions[0].Index).Should(Equal(1))
				Ω(startAuctions[1].Index).Should(Equal(2))
				Ω(startAuctions[2].Index).Should(Equal(3))
			})

			It("does not stop extra ones", func() {
				Consistently(bbs.GetStopLRPInstances).Should(BeEmpty())
			})
		})

		Context("when there are extra instances running for the desired app", func() {
			BeforeEach(func() {
				desiredLRP.Instances = 2
				bbs.Lock()
				bbs.ActualLRPs = []models.ActualLRP{
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "a",
						Index:        0,
						State:        models.ActualLRPStateStarting,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "b",
						Index:        1,
						State:        models.ActualLRPStateStarting,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "c",
						Index:        2,
						State:        models.ActualLRPStateRunning,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "d",
						Index:        3,
						State:        models.ActualLRPStateRunning,
					},
				}
				bbs.Unlock()
			})

			It("doesn't start anything", func() {
				Consistently(bbs.GetLRPStartAuctions).Should(BeEmpty())
			})

			It("stops extra ones", func() {
				Eventually(bbs.GetStopLRPInstances).Should(HaveLen(2))
				stopInstances := bbs.GetStopLRPInstances()

				stopInstance1 := models.StopLRPInstance{
					ProcessGuid:  "the-app-guid-the-app-version",
					Index:        2,
					InstanceGuid: "c",
				}
				stopInstance2 := models.StopLRPInstance{
					ProcessGuid:  "the-app-guid-the-app-version",
					Index:        3,
					InstanceGuid: "d",
				}

				Ω(stopInstances).Should(ContainElement(stopInstance1))
				Ω(stopInstances).Should(ContainElement(stopInstance2))
			})
		})

		Context("when there are duplicate desired instances running for the desired app", func() {
			BeforeEach(func() {
				desiredLRP.Instances = 3
				bbs.Lock()
				bbs.ActualLRPs = []models.ActualLRP{
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "a",
						Index:        0,
						State:        models.ActualLRPStateStarting,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "b",
						Index:        1,
						State:        models.ActualLRPStateStarting,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "c",
						Index:        1,
						State:        models.ActualLRPStateStarting,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "d",
						Index:        2,
						State:        models.ActualLRPStateRunning,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "e",
						Index:        2,
						State:        models.ActualLRPStateRunning,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "f",
						Index:        3,
						State:        models.ActualLRPStateRunning,
					},
					{
						ProcessGuid:  "the-app-guid-the-app-version",
						InstanceGuid: "g",
						Index:        3,
						State:        models.ActualLRPStateRunning,
					},
				}
				bbs.Unlock()
			})

			It("doesn't start anything", func() {
				Consistently(bbs.GetLRPStartAuctions).Should(BeEmpty())
			})

			It("holds stop auctions for the desired duplicates", func() {
				Eventually(bbs.GetLRPStopAuctions).Should(HaveLen(2))
				stopAuctions := bbs.GetLRPStopAuctions()

				Ω(stopAuctions).Should(ContainElement(models.LRPStopAuction{
					ProcessGuid: "the-app-guid-the-app-version",
					Index:       1,
				}))

				Ω(stopAuctions).Should(ContainElement(models.LRPStopAuction{
					ProcessGuid: "the-app-guid-the-app-version",
					Index:       2,
				}))
			})

			It("stops extra ones", func() {
				Eventually(bbs.GetStopLRPInstances).Should(HaveLen(2))
				stopInstances := bbs.GetStopLRPInstances()

				stopInstance1 := models.StopLRPInstance{
					ProcessGuid:  "the-app-guid-the-app-version",
					Index:        3,
					InstanceGuid: "f",
				}
				stopInstance2 := models.StopLRPInstance{
					ProcessGuid:  "the-app-guid-the-app-version",
					Index:        3,
					InstanceGuid: "g",
				}

				Ω(stopInstances).Should(ContainElement(stopInstance1))
				Ω(stopInstances).Should(ContainElement(stopInstance2))
			})
		})
	})

	Describe("when a desired LRP is deleted", func() {
		JustBeforeEach(func() {
			bbs.DesiredLRPChangeChan <- models.DesiredLRPChange{
				Before: &desiredLRP,
				After:  nil,
			}
		})

		BeforeEach(func() {
			bbs.Lock()
			bbs.ActualLRPs = []models.ActualLRP{
				{
					ProcessGuid:  "the-app-guid-the-app-version",
					InstanceGuid: "a",
					Index:        0,
					State:        models.ActualLRPStateStarting,
				},
			}
			bbs.Unlock()
		})

		It("doesn't start anything", func() {
			Consistently(bbs.GetLRPStartAuctions).Should(BeEmpty())
		})

		It("stops all instances", func() {
			Eventually(bbs.GetStopLRPInstances).Should(HaveLen(1))
			stopInstances := bbs.GetStopLRPInstances()

			stopInstance := models.StopLRPInstance{
				ProcessGuid:  "the-app-guid-the-app-version",
				Index:        0,
				InstanceGuid: "a",
			}

			Ω(stopInstances).Should(ContainElement(stopInstance))
		})
	})
})
