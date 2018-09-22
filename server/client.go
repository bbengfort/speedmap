package server

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/bbengfort/speedmap/server/pb"
	"google.golang.org/grpc"
)

// DefaultTimeout for making client requests to the speedmap server
const DefaultTimeout = 30 * time.Second

// Client is a helper struct to connect to the speedmap server and make requests.
type Client struct {
	identity string
	conn     *grpc.ClientConn
	client   pb.KVClient
}

// NewClient creates a new speedmap server client and returns it
func NewClient(identity string) *Client {
	client := &Client{identity: identity}

	if client.identity == "" {
		hostname, _ := os.Hostname()
		if hostname != "" {
			client.identity = fmt.Sprintf("%s-%04X", hostname, rand.Intn(0x10000))
		} else {
			client.identity = fmt.Sprintf("%04X-%04X", rand.Intn(0x10000), rand.Intn(0x10000))
		}
	}

	return client
}

//===========================================================================
// Connection Handling
//===========================================================================

// Connect to the speedmap server and prepare to make requests
func (c *Client) Connect(addr string) (err error) {
	// Close the connection if one is already open.
	c.Close()

	// Connect to the specified address
	if c.conn, err = grpc.Dial(addr, grpc.WithInsecure(), grpc.WithTimeout(DefaultTimeout)); err != nil {
		return err
	}

	// Create the grpc client
	c.client = pb.NewKVClient(c.conn)
	return nil
}

// Close the connection to the speedmap server
func (c *Client) Close() (err error) {
	// Ensure valid state after close
	defer func() {
		c.conn = nil
		c.client = nil
	}()

	if c.conn != nil {
		if err = c.conn.Close(); err != nil {
			return err
		}
	}

	return nil
}

//===========================================================================
// Client Requests
//===========================================================================

// Get performs a request to the speedmap server for the specified key.
func (c *Client) Get(key string) (*pb.ClientReply, error) {
	// Ensure that we're connected
	if c.client == nil {
		return nil, errors.New("not connected to speedmap server")
	}

	// Create the request
	req := &pb.GetRequest{
		Identity: c.identity,
		Key:      key,
	}

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return c.client.Get(ctx, req)
}

// Put performs a request to the speedmap server for the specified key and value.
func (c *Client) Put(key string, value []byte) (*pb.ClientReply, error) {
	// Ensure that we're connected
	if c.client == nil {
		return nil, errors.New("not connected to speedmap server")
	}

	// Create the request
	req := &pb.PutRequest{
		Identity: c.identity,
		Key:      key,
		Value:    value,
	}

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return c.client.Put(ctx, req)
}

// Del performs a request to the speedmap server for the specified key.
func (c *Client) Del(key string, force bool) (*pb.ClientReply, error) {
	// Ensure that we're connected
	if c.client == nil {
		return nil, errors.New("not connected to speedmap server")
	}

	// Create the request
	req := &pb.DelRequest{
		Identity: c.identity,
		Key:      key,
		Force:    force,
	}

	// Create the context
	ctx, cancel := context.WithTimeout(context.Background(), DefaultTimeout)
	defer cancel()

	return c.client.Del(ctx, req)
}
