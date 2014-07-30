package handler

import (
	"errors"
	"os"
	"sync"

	"github.com/cloudfoundry-incubator/app-manager/stop_message_builder"
	"github.com/cloudfoundry-incubator/delta_force/delta_force"
	Bbs "github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-golang/lager"
)

var ErrNoHealthCheckDefined = errors.New("no health check defined for stack")

type LRPreProcessor interface {
	PreProcess(lrp models.DesiredLRP, instanceIndex int, instanceGuid string) (models.DesiredLRP, error)
}

type Handler struct {
	bbs                Bbs.AppManagerBBS
	lrPreProcessor     LRPreProcessor
	numAZs int
	stopMessageBuilder *stop_message_builder.StopMessageBuilder
	logger             lager.Logger
}

func NewHandler(
	bbs Bbs.AppManagerBBS,
	lrPreProcessor LRPreProcessor,
	numAZs int,
	stopMessageBuilder *stop_message_builder.StopMessageBuilder,
	logger lager.Logger,
) Handler {
	handlerLogger := logger.Session("handler")
	return Handler{
		bbs:                bbs,
		lrPreProcessor:     lrPreProcessor,
		numAZs: numAZs,
		stopMessageBuilder: stopMessageBuilder,
		logger:             handlerLogger,
	}
}

func (h Handler) Run(signals <-chan os.Signal, ready chan<- struct{}) error {
	wg := new(sync.WaitGroup)
	desiredChangeChan, stopChan, errChan := h.bbs.WatchForDesiredLRPChanges()

	close(ready)

	for {
		if desiredChangeChan == nil {
			desiredChangeChan, stopChan, errChan = h.bbs.WatchForDesiredLRPChanges()
		}

		select {
		case desiredChange, ok := <-desiredChangeChan:
			if ok {
				wg.Add(1)
				go func() {
					defer wg.Done()
					h.processDesiredChange(desiredChange)
				}()
			} else {
				h.logger.Error("watch-closed", nil)
				desiredChangeChan = nil
			}

		case err, ok := <-errChan:
			if ok {
				h.logger.Error("watch-error", err)
			}
			desiredChangeChan = nil

		case <-signals:
			h.logger.Info("shutting-down")
			close(stopChan)
			wg.Wait()
			h.logger.Info("shut-down")
			return nil
		}
	}

	return nil
}

func (h Handler) processDesiredChange(desiredChange models.DesiredLRPChange) {
	var desiredLRP models.DesiredLRP
	var desiredInstances int

	changeLogger := h.logger.Session("desired-lrp-change")

	if desiredChange.After == nil {
		desiredLRP = *desiredChange.Before
		desiredInstances = 0
	} else {
		desiredLRP = *desiredChange.After
		desiredInstances = desiredLRP.Instances
	}

	actualInstances, instanceGuidToActual, err := h.actualsForProcessGuid(desiredLRP.ProcessGuid)
	if err != nil {
		changeLogger.Error("fetch-actuals-failed", err, lager.Data{"desired-app-message": desiredLRP})
		return
	}

	delta := delta_force.Reconcile(desiredInstances, actualInstances)

	for _, lrpIndex := range delta.IndicesToStart {
		changeLogger.Info("request-start", lager.Data{
			"desired-app-message": desiredLRP,
			"index":               lrpIndex,
		})

		instanceGuid, err := uuid.NewV4()
		if err != nil {
			changeLogger.Error("generating-instance-guid-failed", err)
			return
		}

		preprocessedLRP, err := h.lrPreProcessor.PreProcess(desiredLRP, lrpIndex, instanceGuid.String())
		if err != nil {
			changeLogger.Error("failed-to-preprocess-lrp", err)
			return
		}

		startMessage := models.LRPStartAuction{
			DesiredLRP: preprocessedLRP,

			NumAZs: h.numAZs,
			Index:        lrpIndex,
			InstanceGuid: instanceGuid.String(),
		}

		err = h.bbs.RequestLRPStartAuction(startMessage)

		if err != nil {
			changeLogger.Error("request-start-auction-failed", err, lager.Data{
				"desired-app-message": desiredLRP,
			"index":               lrpIndex,
			})

		}
	}

	for _, guidToStop := range delta.GuidsToStop {
		changeLogger.Info("request-stop-instance", lager.Data{
			"desired-app-message": desiredLRP,
			"stop-instance-guid":  guidToStop,
		})

		actualToStop := instanceGuidToActual[guidToStop]

		err = h.bbs.RequestStopLRPInstance(models.StopLRPInstance{
			ProcessGuid:  actualToStop.ProcessGuid,
			InstanceGuid: actualToStop.InstanceGuid,
			Index:        actualToStop.Index,
		})

		if err != nil {
			changeLogger.Error("request-stop-instance-failed", err, lager.Data{
				"desired-app-message": desiredLRP,
				"stop-instance-guid":  guidToStop,
			})
		}
	}

	for _, indexToStopAllButOne := range delta.IndicesToStopAllButOne {
		changeLogger.Info("request-stop-auction", lager.Data{
			"desired-app-message":  desiredLRP,
			"stop-duplicate-index": indexToStopAllButOne,
		})

		stopMessage := h.stopMessageBuilder.Build(desiredLRP, indexToStopAllButOne)
		err = h.bbs.RequestLRPStopAuction(stopMessage)

		if err != nil {
			changeLogger.Error("request-stop-auction-failed", err, lager.Data{
				"desired-app-message":  desiredLRP,
				"stop-duplicate-index": indexToStopAllButOne,
			})
		}
	}
}

func (h Handler) actualsForProcessGuid(lrpGuid string) (delta_force.ActualInstances, map[string]models.ActualLRP, error) {
	actualInstances := delta_force.ActualInstances{}
	actualLRPs, err := h.bbs.GetActualLRPsByProcessGuid(lrpGuid)
	instanceGuidToActual := map[string]models.ActualLRP{}

	if err != nil {
		return actualInstances, instanceGuidToActual, err
	}

	for _, actualLRP := range actualLRPs {
		actualInstances = append(actualInstances, delta_force.ActualInstance{actualLRP.Index, actualLRP.InstanceGuid})
		instanceGuidToActual[actualLRP.InstanceGuid] = actualLRP
	}

	return actualInstances, instanceGuidToActual, err
}
