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
				xdsData.SetResources(nil, logger)
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

func setResourceSlice(xdsData *model.XDSData, rt model.ResourceType, nodes []types.Any, logger zerolog.Logger) {
	var collection []interface{}
	switch rt {
	case model.Cluster:
		for _, node := range nodes {
			if resource, ok := appendCluster(node, logger); ok {
				collection = append(collection, resource)
			}
		}
	case model.ClusterLoadAssignment:
		for _, node := range nodes {
			if resource, ok := appendCLA(node, logger); ok {
				collection = append(collection, resource)
			}
		}
	case model.RouteConfiguration:
		for _, node := range nodes {
			if resource, ok := appendRouteConfiguration(node, logger); ok {
				collection = append(collection, resource)
			}
		}
	case model.Listener:
		for _, node := range nodes {
			if resource, ok := appendListener(node, logger); ok {
				collection = append(collection, resource)
			}
		}
	case model.Secret:
		for _, node := range nodes {
			if resource, ok := appendSecret(node, logger); ok {
				collection = append(collection, resource)
			}
		}
	}
	xdsData.SetResources(collection, logger)
}

func appendCluster(node types.Any, logger zerolog.Logger) (v2.Cluster, bool) {
	var resource v2.Cluster
	if err := types.UnmarshalAny(&node, &resource); err != nil {
		logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
		return v2.Cluster{}, false
	}
	return resource, true
}

func appendCLA(node types.Any, logger zerolog.Logger) (v2.ClusterLoadAssignment, bool) {
	var resource v2.ClusterLoadAssignment
	if err := types.UnmarshalAny(&node, &resource); err != nil {
		logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
		return v2.ClusterLoadAssignment{}, false
	}
	return resource, true
}

func appendRouteConfiguration(node types.Any, logger zerolog.Logger) (v2.RouteConfiguration, bool) {
	var resource v2.RouteConfiguration
	if err := types.UnmarshalAny(&node, &resource); err != nil {
		logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
		return v2.RouteConfiguration{}, false
	}
	return resource, true
}

func appendListener(node types.Any, logger zerolog.Logger) (v2.Listener, bool) {
	var resource v2.Listener
	if err := types.UnmarshalAny(&node, &resource); err != nil {
		logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
		return v2.Listener{}, false
	}
	return resource, true
}

func appendSecret(node types.Any, logger zerolog.Logger) (auth.Secret, bool) {
	var resource auth.Secret
	if err := types.UnmarshalAny(&node, &resource); err != nil {
		logger.Error().AnErr("types.UnmarshalAny", err).Msg("Could not unmarshal proto")
		return auth.Secret{}, false
	}
	return resource, true
}
