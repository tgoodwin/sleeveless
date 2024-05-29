package client

import (
	"context"
	"fmt"

	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/go-logr/logr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
)

var log = logf.Log.WithName("sleeveless")

func createFixedLengthHash() string {
	// Get the current time
	currentTime := time.Now()

	// Convert the current time to a byte slice
	timeBytes := []byte(currentTime.String())

	// Hash the byte slice using SHA256
	hash := sha256.Sum256(timeBytes)

	// Convert the hash to a fixed length string
	hashString := hex.EncodeToString(hash[:])

	// Take the first 6 characters of the hash string
	shortHash := hashString[:6]

	return shortHash
}

type Client struct {
	// this syntax is "embedding" the client.Client interface in the Client struct
	// this means that the Client struct will have all the methods of the client.Client interface.
	// below, we will override some of these methods to add our own behavior.
	client.Client

	// used to scope observations to a given Reconcile invocation
	reconcileID string

	logger logr.Logger

	// used for causality propagation -- which objects caused the creation of which objects
	lc LabelContext
}

var _ client.Client = &Client{}

func newClient(wrapped client.Client) client.Client {
	return &Client{
		Client: wrapped,
		logger: log,
	}
}

func Wrap(c client.Client) client.Client {
	return newClient(c)
}

func (c *Client) StartReconcileContext() func() {
	if c.reconcileID != "" {
		panic("concurrent reconcile invocations detected")
	}
	c.reconcileID = createFixedLengthHash()
	return func() {
		c.logger.WithValues(
			"ReconcileID", c.reconcileID,
			"Timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond)),
		).Info("Reconcile context ended")

		c.reconcileID = ""
	}
}

func (c *Client) logObservation(ov ObjectVersion, msg string) {
	c.logger.WithValues(
		"Timestamp", fmt.Sprintf("%d", time.Now().UnixNano()/int64(time.Millisecond)),
		"ReconcileID", c.reconcileID,
		"ObservedObjectKind", fmt.Sprintf("%+v", ov.Kind),
		"ObservedObjectUID", fmt.Sprintf("%+v", ov.Uid),
		"ObservedObjectVersion", fmt.Sprintf("%+v", ov.Version),
		"ObservationTraceID", fmt.Sprintf("%+v", ov.TraceID),
	).Info(msg)
}

func (c *Client) setLabelContext(obj client.Object) {
	labels := obj.GetLabels()
	rootID, ok := labels[TRACEY_ROOT_ID]
	c.lc.SourceObject = string(obj.GetUID())
	if !ok {
		return
	}
	c.lc.RootID = rootID
	if _, ok := labels[TRACEY_PARENT_ID]; !ok {
		c.lc.ParentID = rootID
	}
	if traceID, ok := labels[TRACEY_LABEL_ID]; ok {
		c.lc.TraceID = traceID
	}
	c.logger.WithValues(
		"RootID", c.lc.RootID,
		"ParentID", c.lc.ParentID,
		"TraceID", c.lc.TraceID,
	).Info("Label context set")
}

func (c *Client) propagateLabels(obj client.Object) {
	currLabels := obj.GetLabels()
	out := make(map[string]string)
	for k, v := range currLabels {
		out[k] = v
	}
	out[TRACEY_PARENT_ID] = c.lc.TraceID
	out[TRACEY_LABEL_ID] = createFixedLengthHash()
	c.logger.WithValues(
		"RootID", c.lc.RootID,
		"ParentID", c.lc.ParentID,
		"TraceID", c.lc.TraceID,
		"SourceObject", c.lc.SourceObject,
		"DestObject", obj.GetUID(),
	).Info("Propagating labels")
	pre := obj.GetLabels()
	obj.SetLabels(out)
	post := obj.GetLabels()
	c.logger.WithValues(
		"PreLabel", pre,
		"PostLabel", post,
	).Info("Labels propagated")
}

func (c *Client) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	c.propagateLabels(obj)
	res := c.Client.Create(ctx, obj, opts...)
	c.logObservation(RecordSingle(obj), "CREATE")
	return res

}

func (c *Client) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	// need to record the knowledge snapshot this delete is based on
	c.propagateLabels(obj)
	res := c.Client.Delete(ctx, obj, opts...)
	c.logObservation(RecordSingle(obj), "DELETE")
	return res
}

func (c *Client) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	res := c.Client.Get(ctx, key, obj, opts...)
	// after the read, set the label context for the next operation
	c.setLabelContext(obj)
	// and then record a knowledge snapshot
	c.logObservation(RecordSingle(obj), "GET")
	return res
}

func (c *Client) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	return c.Client.List(ctx, list, opts...)
}

func (c *Client) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	// need to record the knowledge snapshot this update is based on
	c.propagateLabels(obj)
	ov := RecordSingle(obj)
	c.logObservation(ov, "UPDATE")
	res := c.Client.Update(ctx, obj, opts...)
	return res
}

func (c *Client) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	// TODO verify labels propagate correctly under patch
	c.propagateLabels(obj)
	res := c.Client.Patch(ctx, obj, patch, opts...)
	ov := RecordSingle(obj)
	c.logObservation(ov, "PATCH")
	return res
}
