package model

import (
	"encoding/json"
	"fmt"
	"sync"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
	auth "github.com/envoyproxy/go-control-plane/envoy/api/v2/auth"
)

type XDSData struct {
	sync.Mutex
	currentType            ResourceType
	clusters               []v2.Cluster
	clusterLoadAssignments []v2.ClusterLoadAssignment
	routeConfigurations    []v2.RouteConfiguration
	listeners              []v2.Listener
	secrets                []auth.Secret
}

func New() *XDSData {
	return &XDSData{}
}

func (d *XDSData) SetClusters(clusters []v2.Cluster) {
	d.Lock()
	d.currentType = Cluster
	d.clusters = clusters
	d.Unlock()
}

func (d *XDSData) SetCLAs(clusterLoadAssignments []v2.ClusterLoadAssignment) {
	d.Lock()
	d.currentType = ClusterLoadAssignment
	d.clusterLoadAssignments = clusterLoadAssignments
	d.Unlock()
}

func (d *XDSData) SetRouteConfigurations(routeConfigurations []v2.RouteConfiguration) {
	d.Lock()
	d.currentType = RouteConfiguration
	d.routeConfigurations = routeConfigurations
	d.Unlock()
}

func (d *XDSData) SetListeners(listeners []v2.Listener) {
	d.Lock()
	d.currentType = Listener
	d.listeners = listeners
	d.Unlock()
}

func (d *XDSData) SetSecrets(secrets []auth.Secret) {
	d.Lock()
	d.currentType = Secret
	d.secrets = secrets
	d.Unlock()
}

func (d *XDSData) GetLatestResources() ([]byte, error) {
	d.Lock()
	defer d.Unlock()

	switch d.currentType {
	case Cluster:
		collection := d.clusters
		if collection == nil {
			collection = make([]v2.Cluster, 0)
		}
		return json.Marshal(collection)

	case ClusterLoadAssignment:
		collection := d.clusterLoadAssignments
		if collection == nil {
			collection = make([]v2.ClusterLoadAssignment, 0)
		}
		return json.Marshal(collection)

	case RouteConfiguration:
		collection := d.routeConfigurations
		if collection == nil {
			collection = make([]v2.RouteConfiguration, 0)
		}
		return json.Marshal(collection)

	case Listener:
		collection := d.listeners
		if collection == nil {
			collection = make([]v2.Listener, 0)
		}
		return json.Marshal(collection)

	case Secret:
		collection := d.secrets
		if collection == nil {
			collection = make([]auth.Secret, 0)
		}
		return json.Marshal(collection)

	default:
		return nil, fmt.Errorf("Invalid ResourceType %+v", d.currentType)
	}

}
