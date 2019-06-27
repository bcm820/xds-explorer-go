package stream

import (
	"context"
	"reflect"

	"github.com/bcmendoza/xds-explorer/models"
	"github.com/deciphernow/gm-fabric-go/discovery"

	"github.com/rs/zerolog"
)

func Listen(ctx context.Context, requestChan <-chan models.Request, logger zerolog.Logger) {
	var err error
	var session *discovery.XDS
	var request *models.Request

	resourceChan := make(chan discovery.Resource, 1)

RESOURCE_LOOP:
	for {
		select {
		case <-ctx.Done():
			break RESOURCE_LOOP

		case req := <-requestChan:

			// TODO: validate fields prior to making GRPC request
			options := []discovery.Option{
				discovery.WithLocation("localhost:50000"),
				discovery.WithResourceType(request.ResourceType),
				discovery.WithZone(request.Zone),
				discovery.WithCluster(request.Cluster),
				discovery.WithResourceNames(request.ResourceNames),
			}

			// compare current request and only proceed if it's different
			if request != nil && reflect.DeepEqual(&req, request) {
				continue RESOURCE_LOOP
			}
			request = &req

			// close existing session
			// if closing fails, don't replace the current session
			if session != nil {
				if err = session.Close(); err != nil {
					logger.Error().AnErr("session.Close()", err).Msg("XDS server session close error")
					continue RESOURCE_LOOP
				}
				logger.Info().Msg("Close XDS session")
			}

			// create new session
			session, err = discovery.NewXDSSession(options...)
			if err != nil {
				logger.Error().AnErr("discovery.NewXDSSession()", err).Msg("XDS server error")
			} else {
				logger.Info().Msg("New XDS session")
				discovery.DiscoverNodesStream(session, resourceChan)
			}

		case <-resourceChan:
			logger.Info().Msg("Incoming from XDS!")
			// TODO: Update internal model with resources
		}
	}

	if err = session.Close(); err != nil {
		logger.Error().AnErr("session.Close()", err).Msg("XDS server session close error")
	} else {
		logger.Info().Msg("Close XDS session")
	}
}
