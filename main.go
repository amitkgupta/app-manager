package main

import (
	"flag"
	"os"
	"strings"

	"github.com/cloudfoundry-incubator/cf-lager"
	Bbs "github.com/cloudfoundry-incubator/runtime-schema/bbs"
	"github.com/cloudfoundry/gunk/timeprovider"
	"github.com/cloudfoundry/storeadapter/etcdstoreadapter"
	"github.com/cloudfoundry/storeadapter/workerpool"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/ifrit"
	"github.com/tedsuo/ifrit/grouper"
	"github.com/tedsuo/ifrit/sigmon"

	"github.com/cloudfoundry-incubator/app-manager/handler"
	"github.com/cloudfoundry-incubator/app-manager/lrpreprocessor"
	"github.com/cloudfoundry-incubator/app-manager/stop_message_builder"
)

var etcdCluster = flag.String(
	"etcdCluster",
	"http://127.0.0.1:4001",
	"comma-separated list of etcd addresses (http://ip:port)",
)

var numAZs = flag.Int(
	"numAZs",
	-1,
	"total number of AZs on which Executors are running",
)

func main() {
	flag.Parse()

	if *numAZs < 0 {
		panic("needs (non-negative) number of AZs")
	}

	logger := cf_lager.New("app-manager")

	bbs := initializeBbs(logger)

	lrpp := lrpreprocessor.New(bbs)

	stopMessageBuilder := stop_message_builder.New(*numAZs)

	group := grouper.EnvokeGroup(grouper.RunGroup{
		"handler": handler.NewHandler(bbs, lrpp, *numAZs, stopMessageBuilder, logger),
	})

	logger.Info("started")

	monitor := ifrit.Envoke(sigmon.New(group))

	err := <-monitor.Wait()
	if err != nil {
		logger.Error("exited-with-failure", err)
		os.Exit(1)
	}

	logger.Info("exited")
}

func initializeBbs(logger lager.Logger) Bbs.AppManagerBBS {
	etcdAdapter := etcdstoreadapter.NewETCDStoreAdapter(
		strings.Split(*etcdCluster, ","),
		workerpool.NewWorkerPool(10),
	)

	err := etcdAdapter.Connect()
	if err != nil {
		logger.Fatal("failed-to-connect-to-etcd", err)
	}

	return Bbs.NewAppManagerBBS(etcdAdapter, timeprovider.NewTimeProvider(), logger)
}
