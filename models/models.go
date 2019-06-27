package models

// Action to take on the retrieve
type Action int

const (
	Start = iota + 1
	Stop
)

type Request struct {
	ResourceType  string   `json:"resourceType"`
	Zone          string   `json:"zone"`
	Cluster       string   `json:"cluster"`
	ResourceNames []string `json:"resourceNames"`
}

/*

	typePrefix   = "type.googleapis.com/envoy.api.v2."
	EndpointType = typePrefix + "ClusterLoadAssignment"
	ClusterType  = typePrefix + "Cluster"
	RouteType    = typePrefix + "RouteConfiguration"
	ListenerType = typePrefix + "Listener"
	SecretType   = typePrefix + "auth.Secret"

		discovery.WithLocation(xds.serverAddress),
		discovery.WithResourceType(cache.EndpointType),
		discovery.WithZone(request.MeshID),
		discovery.WithCluster(request.ClusterID),
		discovery.WithResourceNames([]string{request.ClusterID})
*/
