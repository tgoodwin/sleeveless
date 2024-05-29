package client

import (
	"github.com/google/uuid"
)

var TRACEY_ROOT_ID = "tracey-uid" // this is set by the webhook only
var TRACEY_LABEL_ID = "discrete.events/trace-id"
var TRACEY_PARENT_ID = "discrete.events/parent-id"

type LabelContext struct {
	RootID       string
	TraceID      string
	ParentID     string
	SourceObject string
}

func propagateLabels2(labels map[string]string) map[string]string {
	out := make(map[string]string)
	for k, v := range labels {
		out[k] = v
	}
	rootID, ok := labels[TRACEY_ROOT_ID]
	if !ok {
		return labels
	}
	// handle the base case where there is no parent yet. The rootID will be the parent to the child
	if _, ok := labels[TRACEY_PARENT_ID]; !ok {
		labels[TRACEY_LABEL_ID] = rootID
	}
	// otherwise, generate a new trace ID and assign the current trace ID as the parent.
	if traceID, ok := labels[TRACEY_LABEL_ID]; ok {
		out[TRACEY_PARENT_ID] = traceID
		out[TRACEY_LABEL_ID] = uuid.New().String()
	}
	return out
}
