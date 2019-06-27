package model

type Request struct {
	ResourceType  string   `json:"resourceType"`
	Node          string   `json:"node"`
	Zone          string   `json:"zone"`
	Cluster       string   `json:"cluster"`
	ResourceNames []string `json:"resourceNames"`
}

/*
	typePrefix   = "type.googleapis.com/envoy.api.v2."
	EndpointType = typePrefix + "ClusterLoadAssignment"

	// TODO
	ClusterType  = typePrefix + "Cluster"
	RouteType    = typePrefix + "RouteConfiguration"
	ListenerType = typePrefix + "Listener"
	SecretType   = typePrefix + "auth.Secret"
*/
