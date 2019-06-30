package model

type ResourceType string

const (
	Cluster               = "Cluster"
	ClusterLoadAssignment = "ClusterLoadAssignment"
	RouteConfiguration    = "RouteConfiguration"
	Listener              = "Listener"
	Secret                = "auth.Secret"
)

type Request struct {
	ResourceType  ResourceType `json:"resourceType"`
	Node          string       `json:"node"`
	Zone          string       `json:"zone"`
	Cluster       string       `json:"cluster"`
	ResourceNames []string     `json:"resourceNames"`
}
