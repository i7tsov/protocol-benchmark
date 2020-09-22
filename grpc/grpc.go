package grpc

import (
	"context"
	"io"
	"time"

	"github.com/i7tsov/protocol-benchmark/pb"
	"github.com/i7tsov/protocol-benchmark/util"
	"google.golang.org/grpc"
)

// Generate ...
func Generate(count int) []pb.Element {
	res := make([]pb.Element, count)
	for i := range res {
		res[i].Name = util.RandString()
		res[i].Class = util.RandString()
		res[i].Subclass = util.RandString()
		res[i].Indicator1 = util.RandInt()
		res[i].Indicator2 = util.RandInt()
	}
	return res
}

// Server ...
type Server struct {
	Elements []pb.Element
}

// Download rpc.
func (s Server) Download(rq *pb.DownloadRequest, srv pb.Server_DownloadServer) error {
	for i := range s.Elements {
		e := pb.Element{
			Name:       s.Elements[i].Name,
			Class:      s.Elements[i].Class,
			Subclass:   s.Elements[i].Subclass,
			Indicator1: s.Elements[i].Indicator1,
			Indicator2: s.Elements[i].Indicator2,
		}
		err := srv.Send(&e)
		if err != nil {
			return err
		}
	}
	return nil
}

// Client ...
type Client struct {
	Conn *grpc.ClientConn
}

// Download ...
func (c *Client) Download() ([]pb.Element, error) {
	cl := pb.NewServerClient(c.Conn)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	stream, err := cl.Download(ctx, &pb.DownloadRequest{})
	if err != nil {
		return nil, err
	}
	elements := make([]pb.Element, 0, 1000)
	for {
		el, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		elements = append(elements, pb.Element{
			Name:       el.Name,
			Class:      el.Class,
			Subclass:   el.Subclass,
			Indicator1: el.Indicator1,
			Indicator2: el.Indicator2,
		})
	}
	return elements, nil
}
