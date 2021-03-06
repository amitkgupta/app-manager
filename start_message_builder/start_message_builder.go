package start_message_builder

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"strings"

	RepRoutes "github.com/cloudfoundry-incubator/rep/routes"
	"github.com/cloudfoundry-incubator/runtime-schema/models"
	SchemaRouter "github.com/cloudfoundry-incubator/runtime-schema/router"
	"github.com/cloudfoundry/gunk/urljoiner"
	"github.com/nu7hatch/gouuid"
	"github.com/pivotal-golang/lager"
	"github.com/tedsuo/rata"
)

var ErrNoCircusDefined = errors.New("no lifecycle binary bundle defined for stack")

type StartMessageBuilder struct {
	repAddrRelativeToExecutor string
	logger                    lager.Logger
	circuses                  map[string]string
}

func New(repAddrRelativeToExecutor string, circuses map[string]string, logger lager.Logger) *StartMessageBuilder {
	return &StartMessageBuilder{
		repAddrRelativeToExecutor: repAddrRelativeToExecutor,
		circuses:                  circuses,
		logger:                    logger,
	}
}

func (b *StartMessageBuilder) Build(desiredLRP models.DesiredLRP, lrpIndex int, fileServerURL string) (models.LRPStartAuction, error) {
	lrpGuid := desiredLRP.ProcessGuid

	buildLogger := b.logger.Session("message-builder")

	instanceGuid, err := uuid.NewV4()
	if err != nil {
		buildLogger.Error("generating-instance-guid-failed", err)
		return models.LRPStartAuction{}, err
	}

	circusURL, err := b.circusDownloadURL(desiredLRP.Stack, fileServerURL)
	if err != nil {
		buildLogger.Error("construct-circus-download-url-failed", err, lager.Data{
			"stack": desiredLRP.Stack,
		})

		return models.LRPStartAuction{}, err
	}

	lrpEnv, err := createLrpEnv(desiredLRP.Environment, lrpGuid, lrpIndex)
	if err != nil {
		buildLogger.Error("constructing-env-failed", err)
		return models.LRPStartAuction{}, err
	}

	var numFiles *uint64
	if desiredLRP.FileDescriptors != 0 {
		numFiles = &desiredLRP.FileDescriptors
	}

	repRequests := rata.NewRequestGenerator(
		"http://"+b.repAddrRelativeToExecutor,
		RepRoutes.Routes,
	)

	healthyHook, err := repRequests.CreateRequest(
		RepRoutes.LRPRunning,
		rata.Params{
			"process_guid":  lrpGuid,
			"index":         fmt.Sprintf("%d", lrpIndex),
			"instance_guid": instanceGuid.String(),
		},
		nil,
	)
	if err != nil {
		return models.LRPStartAuction{}, err
	}

	return models.LRPStartAuction{
		ProcessGuid:  lrpGuid,
		InstanceGuid: instanceGuid.String(),
		State:        models.LRPStartAuctionStatePending,
		Index:        lrpIndex,

		MemoryMB: desiredLRP.MemoryMB,
		DiskMB:   desiredLRP.DiskMB,

		Ports: []models.PortMapping{
			{ContainerPort: 8080},
		},

		Stack: desiredLRP.Stack,
		Log: models.LogConfig{
			Guid:       desiredLRP.LogGuid,
			SourceName: "App",
			Index:      &lrpIndex,
		},
		Actions: []models.ExecutorAction{
			{
				Action: models.DownloadAction{
					From:    circusURL.String(),
					To:      "/tmp/circus",
					Extract: true,
				},
			},
			{
				Action: models.DownloadAction{
					From:     desiredLRP.Source,
					To:       ".",
					Extract:  true,
					CacheKey: fmt.Sprintf("droplets-%s", lrpGuid),
				},
			},
			models.Parallel(
				models.ExecutorAction{
					models.RunAction{
						Path:    "/tmp/circus/soldier",
						Args:    append([]string{"/app"}, strings.Split(desiredLRP.StartCommand, " ")...),
						Env:     lrpEnv,
						Timeout: 0,
						ResourceLimits: models.ResourceLimits{
							Nofile: numFiles,
						},
					},
				},
				models.ExecutorAction{
					models.MonitorAction{
						Action: models.ExecutorAction{
							models.RunAction{
								Path: "/tmp/circus/spy",
								Args: []string{"-addr=:8080"},
							},
						},
						HealthyThreshold:   1,
						UnhealthyThreshold: 1,
						HealthyHook: models.HealthRequest{
							Method: healthyHook.Method,
							URL:    healthyHook.URL.String(),
						},
					},
				},
			),
		},
	}, nil
}

func (b StartMessageBuilder) circusDownloadURL(stack string, fileServerURL string) (*url.URL, error) {
	checkPath, ok := b.circuses[stack]
	if !ok {
		return nil, ErrNoCircusDefined
	}

	staticRoute, ok := SchemaRouter.NewFileServerRoutes().RouteForHandler(SchemaRouter.FS_STATIC)
	if !ok {
		return nil, errors.New("couldn't generate the download path for the bundle of app lifecycle binaries")
	}

	urlString := urljoiner.Join(fileServerURL, staticRoute.Path, checkPath)

	url, err := url.ParseRequestURI(urlString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse download URL for app lifecycle binary bundle: %s", err)
	}

	return url, nil
}

func createLrpEnv(env []models.EnvironmentVariable, lrpGuid string, lrpIndex int) ([]models.EnvironmentVariable, error) {
	env = append(env, models.EnvironmentVariable{Name: "PORT", Value: "8080"})
	env = append(env, models.EnvironmentVariable{Name: "VCAP_APP_PORT", Value: "8080"})
	env = append(env, models.EnvironmentVariable{Name: "VCAP_APP_HOST", Value: "0.0.0.0"})

	vcapAppEnv := map[string]interface{}{}
	vcapAppEnvIndex := -1
	for i, envVar := range env {
		if envVar.Name == "VCAP_APPLICATION" {
			vcapAppEnvIndex = i
			err := json.Unmarshal([]byte(envVar.Value), &vcapAppEnv)
			if err != nil {
				return env, err
			}
		}
	}

	if vcapAppEnvIndex == -1 {
		return env, nil
	}

	vcapAppEnv["port"] = 8080
	vcapAppEnv["host"] = "0.0.0.0"
	vcapAppEnv["instance_id"] = lrpGuid
	vcapAppEnv["instance_index"] = lrpIndex

	lrpEnv, err := json.Marshal(vcapAppEnv)
	if err != nil {
		return env, err
	}

	env[vcapAppEnvIndex] = models.EnvironmentVariable{Name: "VCAP_APPLICATION", Value: string(lrpEnv)}
	return env, nil
}
