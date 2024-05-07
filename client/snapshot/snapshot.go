package snapshot

import (
	"time"

	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ObjectVersion struct {
	Uid     string `json:"uid"`
	Version string `json:"version"`
}

func (o ObjectVersion) NewerThan(other ObjectVersion) bool {
	return o.Version > other.Version
}

type VersionSet map[ObjectVersion]struct{}

func (s VersionSet) Contains(other VersionSet) bool {
	for k := range other {
		if _, ok := s[k]; !ok {
			return false
		}
	}
	return true
}

// Intersection can be used to represent common knowledge between two sets of versions.
func (s VersionSet) Intersection(other VersionSet) VersionSet {
	result := make(VersionSet)
	for k := range s {
		if _, ok := other[k]; ok {
			result[k] = struct{}{}
		}
	}
	return result
}

func (s VersionSet) Diff(o VersionSet) VersionSet {
	result := make(VersionSet)
	for k := range s {
		if _, ok := o[k]; !ok {
			result[k] = struct{}{}
		}
	}
	return result
}

// Observation represents the result of a List operation on a Kubernetes resource at a given point in time
type Observation struct {
	*VersionSet
	resourceVersion string
}

func (o Observation) Precedes(other Observation) bool {
	return o.resourceVersion < other.resourceVersion
}

func (o Observation) Contains(other Observation) bool {
	return o.VersionSet.Contains(*other.VersionSet)
}

type key schema.GroupVersionKind

type Snapshot struct {
	Objects         map[key]VersionSet
	Timestamp       time.Time
	ResourceVersion string
}
