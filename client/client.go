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
}

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
	c.reconcileID = createFixedLengthHash()
	return func() {
		c.logger.Info("Reconcile context ended", "ReconcileID", c.reconcileID)
		c.reconcileID = ""
	}
}

func (c *Client) logObservation(ov ObjectVersion, msg string) {
	c.logger.WithValues(
		"Timestamp", time.Now().Format("2006-01-02 15:04:05"),
		"ReconcileID", c.reconcileID,
		"Observation", fmt.Sprintf("%+v", ov),
	).Info(msg)
}

func (c *Client) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	res := c.Client.Create(ctx, obj, opts...)
	ov := RecordSingle(obj)
	c.logObservation(ov, "CREATE")
	return res

}

func (c *Client) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	// need to record the knowledge snapshot this delete is based on
	res := c.Client.Delete(ctx, obj, opts...)
	ov := RecordSingle(obj)
	c.logObservation(ov, "DELETE")
	return res
}

func (c *Client) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	// record a knowledge snapshot
	res := c.Client.Get(ctx, key, obj, opts...)
	ov := RecordSingle(obj)
	c.logObservation(ov, "GET")
	return res
}

func (c *Client) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	fmt.Println("sleeveless LIST")
	return c.Client.List(ctx, list, opts...)
}

func (c *Client) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	// need to record the knowledge snapshot this update is based on
	res := c.Client.Update(ctx, obj, opts...)
	ov := RecordSingle(obj)
	c.logObservation(ov, "UPDATE")
	return res
}

func (c *Client) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	res := c.Client.Patch(ctx, obj, patch, opts...)
	ov := RecordSingle(obj)
	c.logObservation(ov, "PATCH")
	return res
}
