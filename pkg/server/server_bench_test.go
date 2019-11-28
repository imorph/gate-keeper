package server

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/imorph/gate-keeper/pkg/api/gatekeeper"
)

func BenchmarkSimple(b *testing.B) {
	b.StopTimer()
	logger, _ := zap.NewProduction()
	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 999999999, 999999999, 999999999)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			b.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		b.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			b.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
	}
	b.StopTimer()

	s.Stop()

}

func BenchmarkBanned(b *testing.B) {
	b.StopTimer()
	logger, _ := zap.NewProduction()
	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 99, 99, 99)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			b.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		b.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			b.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
	}
	b.StopTimer()

	s.Stop()

}

func BenchmarkBlackListed(b *testing.B) {
	b.StopTimer()
	logger, _ := zap.NewProduction()
	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 999999999, 999999999, 999999999)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			b.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		b.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			b.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)
	okB, err := c.BlackList(context.Background(), &gatekeeper.BlackListRequest{Subnet: "192.168.1.0/24", Isadd: true})
	if err != nil {
		b.Error("Error when calling BlackList:", err)
	}
	if !okB.GetOk() {
		b.Error("Want:", true, "from BlackList, got:", okB.GetOk())
	}
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		_, _ = c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "192.168.1.1"})
	}
	b.StopTimer()

	s.Stop()

}
