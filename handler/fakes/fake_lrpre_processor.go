// This file was generated by counterfeiter
package fakes

import (
	"github.com/cloudfoundry-incubator/app-manager/handler"
	"github.com/cloudfoundry-incubator/runtime-schema/models"

	"sync"
)

type FakeLRPreProcessor struct {
	PreProcessStub        func(lrp models.DesiredLRP, instanceIndex int, instanceGuid string) (models.DesiredLRP, error)
	preProcessMutex       sync.RWMutex
	preProcessArgsForCall []struct {
		lrp           models.DesiredLRP
		instanceIndex int
		instanceGuid  string
	}
	preProcessReturns struct {
		result1 models.DesiredLRP
		result2 error
	}
}

func (fake *FakeLRPreProcessor) PreProcess(lrp models.DesiredLRP, instanceIndex int, instanceGuid string) (models.DesiredLRP, error) {
	fake.preProcessMutex.Lock()
	defer fake.preProcessMutex.Unlock()
	fake.preProcessArgsForCall = append(fake.preProcessArgsForCall, struct {
		lrp           models.DesiredLRP
		instanceIndex int
		instanceGuid  string
	}{lrp, instanceIndex, instanceGuid})
	if fake.PreProcessStub != nil {
		return fake.PreProcessStub(lrp, instanceIndex, instanceGuid)
	} else {
		return fake.preProcessReturns.result1, fake.preProcessReturns.result2
	}
}

func (fake *FakeLRPreProcessor) PreProcessCallCount() int {
	fake.preProcessMutex.RLock()
	defer fake.preProcessMutex.RUnlock()
	return len(fake.preProcessArgsForCall)
}

func (fake *FakeLRPreProcessor) PreProcessArgsForCall(i int) (models.DesiredLRP, int, string) {
	fake.preProcessMutex.RLock()
	defer fake.preProcessMutex.RUnlock()
	return fake.preProcessArgsForCall[i].lrp, fake.preProcessArgsForCall[i].instanceIndex, fake.preProcessArgsForCall[i].instanceGuid
}

func (fake *FakeLRPreProcessor) PreProcessReturns(result1 models.DesiredLRP, result2 error) {
	fake.preProcessReturns = struct {
		result1 models.DesiredLRP
		result2 error
	}{result1, result2}
}

var _ handler.LRPreProcessor = new(FakeLRPreProcessor)
