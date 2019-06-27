package stream

import (
	"context"

	"github.com/bcmendoza/xds-explorer/models"
	"github.com/deciphernow/gm-fabric-go/discovery"

	"github.com/rs/zerolog"
)

func Listen(ctx context.Context, requestChan <-chan models.Request, logger zerolog.Logger) {
	var session *discovery.XDS
	var err error

	resourceChan := make(chan discovery.Resource, 1)

RESOURCE_LOOP:
	for {
		select {
		case req := <-requestChan:
			// close existing session
			if session != nil {
				if err = session.Close(); err != nil {
					logger.Error().AnErr("session.Close()", err).Msg("XDS server session close error")
				} else {
					logger.Info().Msg("Close session with XDS server")
				}
			}

			// create new session
			session, err = discovery.NewXDSSession(
				discovery.WithLocation("localhost:50000"),
				discovery.WithResourceType(req.ResourceType),
				discovery.WithZone(req.Zone),
				discovery.WithCluster(req.Cluster),
				discovery.WithResourceNames(req.ResourceNames),
			)
			if err != nil {
				logger.Error().AnErr("discovery.NewXDSSession()", err).Msg("XDS server error")
			} else {
				logger.Info().Msg("New session with XDS server")
				discovery.DiscoverNodesStream(session, resourceChan)
			}

		case <-resourceChan:
			logger.Info().Msg("Incoming from XDS!")
			// TODO: Update internal model with resources

		case <-ctx.Done():
			logger.Info().Msg("ctx.Done()")
			break RESOURCE_LOOP
		}
	}

	if err = session.Close(); err != nil {
		logger.Error().AnErr("session.Close()", err).Msg("XDS server session close error")
	} else {
		logger.Info().Msg("Close session with XDS server")
	}
}
