package foo

import (
	"context"
	"fmt"
	"net"
	"testing"

	"github.com/influxdata/yarpc"
)

type fooTestServer struct {
}

func (f *fooTestServer) UnaryMethod(ctx context.Context, in *Request) (*Response, error) {
	return &Response{"world"}, nil
}

func (f *fooTestServer) ServerStreamMethod(in *Request, stm Foo_ServerStreamMethodServer) error {
	for i := 0; i < 10; i++ {
		err := stm.Send(&Response{fmt.Sprintf("val %d", i)})
		if err != nil {
			println("server stream error", err.Error())
			return err
		}
	}
	return nil
}

func makeServer(t testing.TB) (l net.Listener, s *yarpc.Server, addr string) {
	var err error
	l, err = net.Listen("tcp", ":4040")
	if err != nil {
		t.Fatal("couldn't start listener", err)
	}

	s = yarpc.NewServer(yarpc.CustomCodec(&fooCodec{}))
	RegisterFooServer(s, &fooTestServer{})

	addr = l.Addr().String()
	return
}

func TestFooClient_UnaryMethod(t *testing.T) {
	l, s, _ := makeServer(t)

	go func() {
		s.Serve(l)
	}()
	defer s.Stop()

	cc, err := yarpc.Dial(":4040", yarpc.WithCodec(&fooCodec{}))
	if err != nil {
		t.Fatal("couldn't dial server", err)
	}

	fs := NewFooClient(cc)
	in := &Request{"hello"}
	val, err := fs.UnaryMethod(context.Background(), in)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	t.Log(*val)
}

func TestFooClient_ServerStreamMethod(t *testing.T) {
	l, s, _ := makeServer(t)

	go func() {
		s.Serve(l)
	}()

	defer s.Stop()

	cc, err := yarpc.Dial(":4040", yarpc.WithCodec(&fooCodec{}))
	if err != nil {
		t.Fatal("couldn't dial server", err)
	}

	fs := NewFooClient(cc)
	in := &Request{"hello"}
	stream, err := fs.ServerStreamMethod(context.Background(), in)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	for {
		val, err := stream.Recv()
		if err != nil {
			t.Log("EOF")
			break
		}

		t.Log(val.Out)
	}
}
