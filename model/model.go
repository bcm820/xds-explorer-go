package model

import (
	"sync"

	v2 "github.com/envoyproxy/go-control-plane/envoy/api/v2"
)

type Resources struct {
	sync.Mutex
	clusterLoadAssignments []v2.ClusterLoadAssignment
}

func New() *Resources {
	return &Resources{}
}

func (r *Resources) SetCLAs(clusterLoadAssignments []v2.ClusterLoadAssignment) {
	r.Lock()
	defer r.Unlock()

	r.clusterLoadAssignments = clusterLoadAssignments
}

func (r *Resources) GetCLAs() []v2.ClusterLoadAssignment {
	r.Lock()
	defer r.Unlock()

	var resources []v2.ClusterLoadAssignment
	for _, cla := range r.clusterLoadAssignments {
		resources = append(resources, cla)
	}

	return resources
}
