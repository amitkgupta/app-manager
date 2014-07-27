package stop_message_builder

import "github.com/cloudfoundry-incubator/runtime-schema/models"

type StopMessageBuilder struct {
	numAZs int
}

func New(numAZs int) *StopMessageBuilder {
	return &StopMessageBuilder{
		numAZs: numAZs,
	}
}

func (b *StopMessageBuilder) Build(desiredLRP models.DesiredLRP, indexToStopAllButOne int) models.LRPStopAuction {
	return models.LRPStopAuction{
		ProcessGuid:  desiredLRP.ProcessGuid,
		Index:        indexToStopAllButOne,
		NumInstances: desiredLRP.Instances,
		NumAZs:       b.numAZs,
	}
}
