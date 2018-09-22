package server

import (
	"context"
	"fmt"
	"net"

	"github.com/bbengfort/speedmap"
	"github.com/bbengfort/speedmap/server/pb"
	"google.golang.org/grpc"
)

// Server implements the KVServer interface and is essentially just a wrapper
// around a speedmap key/value Store. Note that the Server performs NO
// synchronization, it simply passes all requests from all clients to the
// speedmap. Because each request is handled in its own go routine, it is
// expected that the wrapped Store is thread-safe.
type Server struct {
	kv speedmap.Store
}

// New creates a new server with the specified key value store.
func New(kv speedmap.Store) *Server {
	return &Server{kv: kv}
}

// Serve the key/value store with the specified store on the specified addr.
func Serve(kv speedmap.Store, addr string) error {
	srv := New(kv)
	return srv.Listen(addr)
}

// Listen for gRPC requests on the specified address and serve each request
// in its own go routine. The server handlers access the speed map.
func (s *Server) Listen(addr string) error {
	sock, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("could not listen on %s", addr)
	}
	defer sock.Close()
	fmt.Printf("serving the %s store on %s\n", s.kv.String(), addr)

	// Initialize and run the gRPC server in its own thread
	srv := grpc.NewServer()
	pb.RegisterKVServer(srv, s)
	return srv.Serve(sock)
}

// Get handles a get request to the speedmap, relying on the speedmap for
// concurrent synchronization of accesses. Note that Get uses GetoOrCreate
// in the speedmap, storing nil as the default value; this means that this
// method will not return a not found error. This decision was made to more
// completely test the misframe implementation of the Store.
func (s *Server) Get(ctx context.Context, in *pb.GetRequest) (*pb.ClientReply, error) {
	val, _ := s.kv.GetOrCreate(in.Key, nil)
	return &pb.ClientReply{
		Success:  true,
		Redirect: "",
		Error:    "",
		Pair: &pb.KVPair{
			Key:   in.Key,
			Value: val,
		},
	}, nil
}

// Put handles a put request to the speedmap, relying on the speedmap for
// concurrent synchronization of accesses.
func (s *Server) Put(ctx context.Context, in *pb.PutRequest) (*pb.ClientReply, error) {
	rep := &pb.ClientReply{Success: true, Redirect: "", Error: "", Pair: nil}

	if err := s.kv.Put(in.Key, in.Value); err != nil {
		rep.Success = false
		rep.Error = err.Error()
	}

	return rep, nil
}

// Del handles a del request to the speedmap, relying on the speedmap for
// concurrent synchronization of accesses.
func (s *Server) Del(ctx context.Context, in *pb.DelRequest) (*pb.ClientReply, error) {
	rep := &pb.ClientReply{Success: true, Redirect: "", Error: "", Pair: nil}

	if err := s.kv.Delete(in.Key); err != nil {
		rep.Success = false
		rep.Error = err.Error()
	}

	return rep, nil
}
