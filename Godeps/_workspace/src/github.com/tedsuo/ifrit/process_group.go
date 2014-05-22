package ifrit

import (
	"errors"
	"fmt"
	"os"
)

func envokeGroup(rGroup RunGroup) Process {
	p := processGroup{}
	count := len(rGroup)
	mChan := make(memberChan, count)

	for name, runner := range rGroup {
		go mChan.envokeMember(name, runner)
	}
	for i := 0; i < count; i++ {
		m := <-mChan
		p[m.name] = m.process
	}
	return p
}

type member struct {
	name    string
	process Process
}

type memberChan chan member

func (mChan memberChan) envokeMember(name string, runner Runner) {
	mChan <- member{name: name, process: Envoke(runner)}
}

type processGroup map[string]Process

func (group processGroup) Signal(signal os.Signal) {
	for _, p := range group {
		p.Signal(signal)
	}
}

func (group processGroup) waitForGroup() error {
	var errMsg string
	for name, p := range group {
		err := <-p.Wait()
		if err != nil {
			errMsg += fmt.Sprintf("%s: %s/n", name, err)
		}
	}

	var err error
	if errMsg != "" {
		err = errors.New(errMsg)
	}
	return err
}

func (group processGroup) Wait() <-chan error {
	errChan := make(chan error, 1)

	go func() {
		errChan <- group.waitForGroup()
	}()

	return errChan
}
