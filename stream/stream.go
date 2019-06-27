package stream

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/bcmendoza/xds-explorer/model"
	"github.com/deciphernow/gm-fabric-go/discovery"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	"github.com/gogo/protobuf/types"
	"github.com/rs/zerolog"
)

func Listen(ctx context.Context, requestChan <-chan model.Request, resources *model.Resources, logger zerolog.Logger) {
	var err error
	var session *discovery.XDS
	var request *model.Request

	resourceChan := make(chan discovery.Resource, 1)

RESOURCE_LOOP:
	for {
		select {
		case <-ctx.Done():
			break RESOURCE_LOOP

		case req := <-requestChan:

			// compare current request and only proceed if it's different
			if request != nil && reflect.DeepEqual(&req, request) {
				continue RESOURCE_LOOP
			}
			request = &req

			// TODO: determine if fields are missing
			const typePrefix = "type.googleapis.com/envoy.api.v2."
			options := []discovery.Option{
				discovery.WithLocation("localhost:50000"),
				discovery.WithResourceType(typePrefix + req.ResourceType),
				discovery.WithNode(req.Node),
				discovery.WithZone(req.Zone),
				discovery.WithCluster(req.Cluster),
				discovery.WithResourceNames(req.ResourceNames),
			}

			// close existing session
			// if closing fails, don't replace the current session
			if session != nil {
				if err = session.Close(); err != nil {
					logger.Error().AnErr("session.Close()", err).Msg("XDS server session close error")
					continue RESOURCE_LOOP
				}
				logger.Info().Msg("Close XDS session")
			}

			// clear state and create new session
			session, err = discovery.NewXDSSession(options...)
			if err != nil {
				logger.Error().AnErr("discovery.NewXDSSession()", err).Msg("XDS server error")
			} else {
				resources.SetCLAs(make([]v2.ClusterLoadAssignment, 0))
				discovery.DiscoverNodesStream(session, resourceChan)
				logger.Info().
					Str("ResourceType", request.ResourceType).
					Str("Node", request.Node).
					Str("Zone", request.Zone).
					Str("Cluster", request.Cluster).
					Str("ResourceNames", strings.Join(request.ResourceNames, ", ")).
					Msg("Open XDS stream with Aggregated Discovery Service")
			}

		case res := <-resourceChan:
			logger.Info().Msg(fmt.Sprintf("Receive %d nodes from XDS stream", len(res.Nodes)))

			var collection []v2.ClusterLoadAssignment

		NODE_LOOP:
			for _, node := range res.Nodes {
				var resource v2.ClusterLoadAssignment
				if err := types.UnmarshalAny(&node, &resource); err != nil {
					logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
					continue NODE_LOOP
				}
				collection = append(collection, resource)
			}
			resources.SetCLAs(collection)
		}
	}

	if session != nil {
		if err = session.Close(); err != nil {
			logger.Error().AnErr("session.Close()", err).Msg("XDS server session close error")
		} else {
			logger.Info().Msg("Close XDS session")
		}
	}
}
