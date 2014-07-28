// This file was generated by counterfeiter
package fake_bbs

import (
	. "github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/bbs/services_bbs"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	"sync"
	"time"
)

type FakeRepBBS struct {
	MaintainExecutorPresenceStub        func(heartbeatInterval time.Duration, executorPresence models.ExecutorPresence) (services_bbs.Presence, <-chan bool, error)
	maintainExecutorPresenceMutex       sync.RWMutex
	maintainExecutorPresenceArgsForCall []struct {
		arg1 time.Duration
		arg2 models.ExecutorPresence
	}
	maintainExecutorPresenceReturns struct {
		result1 services_bbs.Presence
		result2 <-chan bool
		result3 error
	}
	WatchForDesiredTaskStub        func() (<-chan models.Task, chan<- bool, <-chan error)
	watchForDesiredTaskMutex       sync.RWMutex
	watchForDesiredTaskArgsForCall []struct{}
	watchForDesiredTaskReturns     struct {
		result1 <-chan models.Task
		result2 chan<- bool
		result3 <-chan error
	}
	ClaimTaskStub        func(taskGuid string, executorID string) error
	claimTaskMutex       sync.RWMutex
	claimTaskArgsForCall []struct {
		arg1 string
		arg2 string
	}
	claimTaskReturns struct {
		result1 error
	}
	StartTaskStub        func(taskGuid string, executorID string, containerHandle string) error
	startTaskMutex       sync.RWMutex
	startTaskArgsForCall []struct {
		arg1 string
		arg2 string
		arg3 string
	}
	startTaskReturns struct {
		result1 error
	}
	CompleteTaskStub        func(taskGuid string, failed bool, failureReason string, result string) error
	completeTaskMutex       sync.RWMutex
	completeTaskArgsForCall []struct {
		arg1 string
		arg2 bool
		arg3 string
		arg4 string
	}
	completeTaskReturns struct {
		result1 error
	}
	ReportActualLRPAsStartingStub        func(lrp models.ActualLRP, executorID string) error
	reportActualLRPAsStartingMutex       sync.RWMutex
	reportActualLRPAsStartingArgsForCall []struct {
		arg1 models.ActualLRP
		arg2 string
	}
	reportActualLRPAsStartingReturns struct {
		result1 error
	}
	ReportActualLRPAsRunningStub        func(lrp models.ActualLRP, executorId string) error
	reportActualLRPAsRunningMutex       sync.RWMutex
	reportActualLRPAsRunningArgsForCall []struct {
		arg1 models.ActualLRP
		arg2 string
	}
	reportActualLRPAsRunningReturns struct {
		result1 error
	}
	RemoveActualLRPStub        func(lrp models.ActualLRP) error
	removeActualLRPMutex       sync.RWMutex
	removeActualLRPArgsForCall []struct {
		arg1 models.ActualLRP
	}
	removeActualLRPReturns struct {
		result1 error
	}
	WatchForStopLRPInstanceStub        func() (<-chan models.StopLRPInstance, chan<- bool, <-chan error)
	watchForStopLRPInstanceMutex       sync.RWMutex
	watchForStopLRPInstanceArgsForCall []struct{}
	watchForStopLRPInstanceReturns     struct {
		result1 <-chan models.StopLRPInstance
		result2 chan<- bool
		result3 <-chan error
	}
	ResolveStopLRPInstanceStub        func(stopInstance models.StopLRPInstance) error
	resolveStopLRPInstanceMutex       sync.RWMutex
	resolveStopLRPInstanceArgsForCall []struct {
		arg1 models.StopLRPInstance
	}
	resolveStopLRPInstanceReturns struct {
		result1 error
	}
}

func (fake *FakeRepBBS) MaintainExecutorPresence(arg1 time.Duration, arg2 models.ExecutorPresence) (services_bbs.Presence, <-chan bool, error) {
	fake.maintainExecutorPresenceMutex.Lock()
	defer fake.maintainExecutorPresenceMutex.Unlock()
	fake.maintainExecutorPresenceArgsForCall = append(fake.maintainExecutorPresenceArgsForCall, struct {
		arg1 time.Duration
		arg2 models.ExecutorPresence
	}{arg1, arg2})
	if fake.MaintainExecutorPresenceStub != nil {
		return fake.MaintainExecutorPresenceStub(arg1, arg2)
	} else {
		return fake.maintainExecutorPresenceReturns.result1, fake.maintainExecutorPresenceReturns.result2, fake.maintainExecutorPresenceReturns.result3
	}
}

func (fake *FakeRepBBS) MaintainExecutorPresenceCallCount() int {
	fake.maintainExecutorPresenceMutex.RLock()
	defer fake.maintainExecutorPresenceMutex.RUnlock()
	return len(fake.maintainExecutorPresenceArgsForCall)
}

func (fake *FakeRepBBS) MaintainExecutorPresenceArgsForCall(i int) (time.Duration, models.ExecutorPresence) {
	fake.maintainExecutorPresenceMutex.RLock()
	defer fake.maintainExecutorPresenceMutex.RUnlock()
	return fake.maintainExecutorPresenceArgsForCall[i].arg1, fake.maintainExecutorPresenceArgsForCall[i].arg2
}

func (fake *FakeRepBBS) MaintainExecutorPresenceReturns(result1 services_bbs.Presence, result2 <-chan bool, result3 error) {
	fake.maintainExecutorPresenceReturns = struct {
		result1 services_bbs.Presence
		result2 <-chan bool
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeRepBBS) WatchForDesiredTask() (<-chan models.Task, chan<- bool, <-chan error) {
	fake.watchForDesiredTaskMutex.Lock()
	defer fake.watchForDesiredTaskMutex.Unlock()
	fake.watchForDesiredTaskArgsForCall = append(fake.watchForDesiredTaskArgsForCall, struct{}{})
	if fake.WatchForDesiredTaskStub != nil {
		return fake.WatchForDesiredTaskStub()
	} else {
		return fake.watchForDesiredTaskReturns.result1, fake.watchForDesiredTaskReturns.result2, fake.watchForDesiredTaskReturns.result3
	}
}

func (fake *FakeRepBBS) WatchForDesiredTaskCallCount() int {
	fake.watchForDesiredTaskMutex.RLock()
	defer fake.watchForDesiredTaskMutex.RUnlock()
	return len(fake.watchForDesiredTaskArgsForCall)
}

func (fake *FakeRepBBS) WatchForDesiredTaskReturns(result1 <-chan models.Task, result2 chan<- bool, result3 <-chan error) {
	fake.watchForDesiredTaskReturns = struct {
		result1 <-chan models.Task
		result2 chan<- bool
		result3 <-chan error
	}{result1, result2, result3}
}

func (fake *FakeRepBBS) ClaimTask(arg1 string, arg2 string) error {
	fake.claimTaskMutex.Lock()
	defer fake.claimTaskMutex.Unlock()
	fake.claimTaskArgsForCall = append(fake.claimTaskArgsForCall, struct {
		arg1 string
		arg2 string
	}{arg1, arg2})
	if fake.ClaimTaskStub != nil {
		return fake.ClaimTaskStub(arg1, arg2)
	} else {
		return fake.claimTaskReturns.result1
	}
}

func (fake *FakeRepBBS) ClaimTaskCallCount() int {
	fake.claimTaskMutex.RLock()
	defer fake.claimTaskMutex.RUnlock()
	return len(fake.claimTaskArgsForCall)
}

func (fake *FakeRepBBS) ClaimTaskArgsForCall(i int) (string, string) {
	fake.claimTaskMutex.RLock()
	defer fake.claimTaskMutex.RUnlock()
	return fake.claimTaskArgsForCall[i].arg1, fake.claimTaskArgsForCall[i].arg2
}

func (fake *FakeRepBBS) ClaimTaskReturns(result1 error) {
	fake.claimTaskReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) StartTask(arg1 string, arg2 string, arg3 string) error {
	fake.startTaskMutex.Lock()
	defer fake.startTaskMutex.Unlock()
	fake.startTaskArgsForCall = append(fake.startTaskArgsForCall, struct {
		arg1 string
		arg2 string
		arg3 string
	}{arg1, arg2, arg3})
	if fake.StartTaskStub != nil {
		return fake.StartTaskStub(arg1, arg2, arg3)
	} else {
		return fake.startTaskReturns.result1
	}
}

func (fake *FakeRepBBS) StartTaskCallCount() int {
	fake.startTaskMutex.RLock()
	defer fake.startTaskMutex.RUnlock()
	return len(fake.startTaskArgsForCall)
}

func (fake *FakeRepBBS) StartTaskArgsForCall(i int) (string, string, string) {
	fake.startTaskMutex.RLock()
	defer fake.startTaskMutex.RUnlock()
	return fake.startTaskArgsForCall[i].arg1, fake.startTaskArgsForCall[i].arg2, fake.startTaskArgsForCall[i].arg3
}

func (fake *FakeRepBBS) StartTaskReturns(result1 error) {
	fake.startTaskReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) CompleteTask(arg1 string, arg2 bool, arg3 string, arg4 string) error {
	fake.completeTaskMutex.Lock()
	defer fake.completeTaskMutex.Unlock()
	fake.completeTaskArgsForCall = append(fake.completeTaskArgsForCall, struct {
		arg1 string
		arg2 bool
		arg3 string
		arg4 string
	}{arg1, arg2, arg3, arg4})
	if fake.CompleteTaskStub != nil {
		return fake.CompleteTaskStub(arg1, arg2, arg3, arg4)
	} else {
		return fake.completeTaskReturns.result1
	}
}

func (fake *FakeRepBBS) CompleteTaskCallCount() int {
	fake.completeTaskMutex.RLock()
	defer fake.completeTaskMutex.RUnlock()
	return len(fake.completeTaskArgsForCall)
}

func (fake *FakeRepBBS) CompleteTaskArgsForCall(i int) (string, bool, string, string) {
	fake.completeTaskMutex.RLock()
	defer fake.completeTaskMutex.RUnlock()
	return fake.completeTaskArgsForCall[i].arg1, fake.completeTaskArgsForCall[i].arg2, fake.completeTaskArgsForCall[i].arg3, fake.completeTaskArgsForCall[i].arg4
}

func (fake *FakeRepBBS) CompleteTaskReturns(result1 error) {
	fake.completeTaskReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) ReportActualLRPAsStarting(arg1 models.ActualLRP, arg2 string) error {
	fake.reportActualLRPAsStartingMutex.Lock()
	defer fake.reportActualLRPAsStartingMutex.Unlock()
	fake.reportActualLRPAsStartingArgsForCall = append(fake.reportActualLRPAsStartingArgsForCall, struct {
		arg1 models.ActualLRP
		arg2 string
	}{arg1, arg2})
	if fake.ReportActualLRPAsStartingStub != nil {
		return fake.ReportActualLRPAsStartingStub(arg1, arg2)
	} else {
		return fake.reportActualLRPAsStartingReturns.result1
	}
}

func (fake *FakeRepBBS) ReportActualLRPAsStartingCallCount() int {
	fake.reportActualLRPAsStartingMutex.RLock()
	defer fake.reportActualLRPAsStartingMutex.RUnlock()
	return len(fake.reportActualLRPAsStartingArgsForCall)
}

func (fake *FakeRepBBS) ReportActualLRPAsStartingArgsForCall(i int) (models.ActualLRP, string) {
	fake.reportActualLRPAsStartingMutex.RLock()
	defer fake.reportActualLRPAsStartingMutex.RUnlock()
	return fake.reportActualLRPAsStartingArgsForCall[i].arg1, fake.reportActualLRPAsStartingArgsForCall[i].arg2
}

func (fake *FakeRepBBS) ReportActualLRPAsStartingReturns(result1 error) {
	fake.reportActualLRPAsStartingReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) ReportActualLRPAsRunning(arg1 models.ActualLRP, arg2 string) error {
	fake.reportActualLRPAsRunningMutex.Lock()
	defer fake.reportActualLRPAsRunningMutex.Unlock()
	fake.reportActualLRPAsRunningArgsForCall = append(fake.reportActualLRPAsRunningArgsForCall, struct {
		arg1 models.ActualLRP
		arg2 string
	}{arg1, arg2})
	if fake.ReportActualLRPAsRunningStub != nil {
		return fake.ReportActualLRPAsRunningStub(arg1, arg2)
	} else {
		return fake.reportActualLRPAsRunningReturns.result1
	}
}

func (fake *FakeRepBBS) ReportActualLRPAsRunningCallCount() int {
	fake.reportActualLRPAsRunningMutex.RLock()
	defer fake.reportActualLRPAsRunningMutex.RUnlock()
	return len(fake.reportActualLRPAsRunningArgsForCall)
}

func (fake *FakeRepBBS) ReportActualLRPAsRunningArgsForCall(i int) (models.ActualLRP, string) {
	fake.reportActualLRPAsRunningMutex.RLock()
	defer fake.reportActualLRPAsRunningMutex.RUnlock()
	return fake.reportActualLRPAsRunningArgsForCall[i].arg1, fake.reportActualLRPAsRunningArgsForCall[i].arg2
}

func (fake *FakeRepBBS) ReportActualLRPAsRunningReturns(result1 error) {
	fake.reportActualLRPAsRunningReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) RemoveActualLRP(arg1 models.ActualLRP) error {
	fake.removeActualLRPMutex.Lock()
	defer fake.removeActualLRPMutex.Unlock()
	fake.removeActualLRPArgsForCall = append(fake.removeActualLRPArgsForCall, struct {
		arg1 models.ActualLRP
	}{arg1})
	if fake.RemoveActualLRPStub != nil {
		return fake.RemoveActualLRPStub(arg1)
	} else {
		return fake.removeActualLRPReturns.result1
	}
}

func (fake *FakeRepBBS) RemoveActualLRPCallCount() int {
	fake.removeActualLRPMutex.RLock()
	defer fake.removeActualLRPMutex.RUnlock()
	return len(fake.removeActualLRPArgsForCall)
}

func (fake *FakeRepBBS) RemoveActualLRPArgsForCall(i int) models.ActualLRP {
	fake.removeActualLRPMutex.RLock()
	defer fake.removeActualLRPMutex.RUnlock()
	return fake.removeActualLRPArgsForCall[i].arg1
}

func (fake *FakeRepBBS) RemoveActualLRPReturns(result1 error) {
	fake.removeActualLRPReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeRepBBS) WatchForStopLRPInstance() (<-chan models.StopLRPInstance, chan<- bool, <-chan error) {
	fake.watchForStopLRPInstanceMutex.Lock()
	defer fake.watchForStopLRPInstanceMutex.Unlock()
	fake.watchForStopLRPInstanceArgsForCall = append(fake.watchForStopLRPInstanceArgsForCall, struct{}{})
	if fake.WatchForStopLRPInstanceStub != nil {
		return fake.WatchForStopLRPInstanceStub()
	} else {
		return fake.watchForStopLRPInstanceReturns.result1, fake.watchForStopLRPInstanceReturns.result2, fake.watchForStopLRPInstanceReturns.result3
	}
}

func (fake *FakeRepBBS) WatchForStopLRPInstanceCallCount() int {
	fake.watchForStopLRPInstanceMutex.RLock()
	defer fake.watchForStopLRPInstanceMutex.RUnlock()
	return len(fake.watchForStopLRPInstanceArgsForCall)
}

func (fake *FakeRepBBS) WatchForStopLRPInstanceReturns(result1 <-chan models.StopLRPInstance, result2 chan<- bool, result3 <-chan error) {
	fake.watchForStopLRPInstanceReturns = struct {
		result1 <-chan models.StopLRPInstance
		result2 chan<- bool
		result3 <-chan error
	}{result1, result2, result3}
}

func (fake *FakeRepBBS) ResolveStopLRPInstance(arg1 models.StopLRPInstance) error {
	fake.resolveStopLRPInstanceMutex.Lock()
	defer fake.resolveStopLRPInstanceMutex.Unlock()
	fake.resolveStopLRPInstanceArgsForCall = append(fake.resolveStopLRPInstanceArgsForCall, struct {
		arg1 models.StopLRPInstance
	}{arg1})
	if fake.ResolveStopLRPInstanceStub != nil {
		return fake.ResolveStopLRPInstanceStub(arg1)
	} else {
		return fake.resolveStopLRPInstanceReturns.result1
	}
}

func (fake *FakeRepBBS) ResolveStopLRPInstanceCallCount() int {
	fake.resolveStopLRPInstanceMutex.RLock()
	defer fake.resolveStopLRPInstanceMutex.RUnlock()
	return len(fake.resolveStopLRPInstanceArgsForCall)
}

func (fake *FakeRepBBS) ResolveStopLRPInstanceArgsForCall(i int) models.StopLRPInstance {
	fake.resolveStopLRPInstanceMutex.RLock()
	defer fake.resolveStopLRPInstanceMutex.RUnlock()
	return fake.resolveStopLRPInstanceArgsForCall[i].arg1
}

func (fake *FakeRepBBS) ResolveStopLRPInstanceReturns(result1 error) {
	fake.resolveStopLRPInstanceReturns = struct {
		result1 error
	}{result1}
}

var _ RepBBS = new(FakeRepBBS)
