package stream

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/bcmendoza/xds-explorer/model"
	"github.com/deciphernow/gm-fabric-go/discovery"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"

	"github.com/gogo/protobuf/types"
	"github.com/rs/zerolog"
)

func Listen(
	xdsServerAddress string,
	ctx context.Context,
	requestChan <-chan model.Request,
	xdsData *model.XDSData,
	logger zerolog.Logger,
) {
	var err error
	var stream *discovery.XDS
	var request *model.Request

	resourceChan := make(chan discovery.Resource, 1)

RESOURCE_LOOP:
	for {
		select {
		case <-ctx.Done():
			break RESOURCE_LOOP

		case req := <-requestChan:
			options, ok := resetStream(request, req, xdsServerAddress, stream, logger)
			if !ok {
				continue RESOURCE_LOOP
			}

			// update current request and create new stream
			request = &req
			stream, err = discovery.NewXDSSession(options...)
			if err != nil {
				logger.Error().AnErr("discovery.NewXDSSession()", err).Msg("XDS stream error")
			} else {
				clearResourceSlice(xdsData, req.ResourceType)
				discovery.DiscoverNodesStream(stream, resourceChan)
				logger.Info().
					Str("ResourceType", fmt.Sprintf("%+v", request.ResourceType)).
					Str("Node", request.Node).
					Str("Zone", request.Zone).
					Str("Cluster", request.Cluster).
					Str("ResourceNames", strings.Join(request.ResourceNames, ", ")).
					Msg("Open XDS stream")
			}

		case res := <-resourceChan:
			logger.Info().Msg(fmt.Sprintf("Receive %d nodes from XDS stream", len(res.Nodes)))
			setResourceSlice(xdsData, request.ResourceType, res.Nodes, logger)
		}
	}

	if stream != nil {
		if err = stream.Close(); err != nil {
			logger.Error().AnErr("stream.Close()", err).Msg("XDS stream close error")
		} else {
			logger.Info().Msg("Close XDS stream")
		}
	}
}

func resetStream(request *model.Request, req model.Request, xdsServerAddress string, stream *discovery.XDS, logger zerolog.Logger) ([]discovery.Option, bool) {
	var err error

	// compare current request and only proceed if it's different
	if request != nil && reflect.DeepEqual(&req, request) {
		return nil, false
	}

	const typePrefix = "type.googleapis.com/envoy.api.v2."
	resourceType := fmt.Sprintf("%s%+v", typePrefix, req.ResourceType)
	options := []discovery.Option{
		discovery.WithLocation(xdsServerAddress),
		discovery.WithResourceType(resourceType),
		discovery.WithNode(req.Node),
		discovery.WithZone(req.Zone),
		discovery.WithCluster(req.Cluster),
		discovery.WithResourceNames(req.ResourceNames),
	}

	// close existing stream
	// if closing fails, don't replace the current stream
	if stream != nil {
		if err = stream.Close(); err != nil {
			logger.Error().AnErr("stream.Close()", err).Msg("XDS stream close error")
			return nil, false
		}
		logger.Info().Msg("Close XDS stream")
	}

	return options, true
}

func clearResourceSlice(xdsData *model.XDSData, rt model.ResourceType) {
	switch rt {
	case model.Cluster:
		xdsData.SetClusters(make([]v2.Cluster, 0))
	case model.ClusterLoadAssignment:
		xdsData.SetCLAs(make([]v2.ClusterLoadAssignment, 0))
	case model.RouteConfiguration:
		xdsData.SetRouteConfigurations(make([]v2.RouteConfiguration, 0))
	case model.Listener:
		xdsData.SetListeners(make([]v2.Listener, 0))
	case model.Secret:
		xdsData.SetSecrets(make([]auth.Secret, 0))
	}
}

func setResourceSlice(xdsData *model.XDSData, rt model.ResourceType, nodes []types.Any, logger zerolog.Logger) {
	switch rt {
	case model.Cluster:
		var collection []v2.Cluster
		for _, node := range nodes {
			var resource v2.Cluster
			if err := types.UnmarshalAny(&node, &resource); err != nil {
				logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
				continue
			}
			collection = append(collection, resource)
		}
		xdsData.SetClusters(collection)

	case model.ClusterLoadAssignment:
		var collection []v2.ClusterLoadAssignment
		for _, node := range nodes {
			var resource v2.ClusterLoadAssignment
			if err := types.UnmarshalAny(&node, &resource); err != nil {
				logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
				continue
			}
			collection = append(collection, resource)
		}
		xdsData.SetCLAs(collection)

	case model.RouteConfiguration:
		var collection []v2.RouteConfiguration
		for _, node := range nodes {
			var resource v2.RouteConfiguration
			if err := types.UnmarshalAny(&node, &resource); err != nil {
				logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
				continue
			}
			collection = append(collection, resource)
		}
		xdsData.SetRouteConfigurations(collection)

	case model.Listener:
		var collection []v2.Listener
		for _, node := range nodes {
			var resource v2.Listener
			if err := types.UnmarshalAny(&node, &resource); err != nil {
				logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
				continue
			}
			collection = append(collection, resource)
		}
		xdsData.SetListeners(collection)

	case model.Secret:
		var collection []auth.Secret
		for _, node := range nodes {
			var resource auth.Secret
			if err := types.UnmarshalAny(&node, &resource); err != nil {
				logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
				continue
			}
			collection = append(collection, resource)
		}
		xdsData.SetSecrets(collection)

	}
}
