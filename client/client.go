package main

import (
	"context"
	"fmt"

	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Client struct {
	client.Client
}

func NewClient(wrapped client.Client) client.Client {
	return &Client{
		Client: wrapped,
	}
}

func (c *Client) Create(ctx context.Context, obj client.Object, opts ...client.CreateOption) error {
	// need to record the knowledge snapshot this create is based on
	fmt.Println("sleeveless CREATE")
	return c.Client.Create(ctx, obj, opts...)
}

func (c *Client) Delete(ctx context.Context, obj client.Object, opts ...client.DeleteOption) error {
	// need to record the knowledge snapshot this delete is based on
	fmt.Println("sleeveless DELETE")
	return c.Client.Delete(ctx, obj, opts...)
}

func (c *Client) Get(ctx context.Context, key client.ObjectKey, obj client.Object, opts ...client.GetOption) error {
	// record a knowledge snapshot
	fmt.Println("sleeveless GET")
	return c.Client.Get(ctx, key, obj, opts...)
}

func (c *Client) List(ctx context.Context, list client.ObjectList, opts ...client.ListOption) error {
	fmt.Println("sleeveless LIST")
	return c.Client.List(ctx, list, opts...)
}

func (c *Client) Update(ctx context.Context, obj client.Object, opts ...client.UpdateOption) error {
	fmt.Println("sleeveless UPDATE")
	// need to record the knowledge snapshot this update is based on
	return c.Client.Update(ctx, obj, opts...)
}

func (c *Client) Patch(ctx context.Context, obj client.Object, patch client.Patch, opts ...client.PatchOption) error {
	return c.Client.Patch(ctx, obj, patch, opts...)
}

func (c *Client) snapshot() {
	// record a knowledge snapshot
	fmt.Println("TODO")
}
