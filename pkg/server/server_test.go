package server

import (
	"testing"
	"time"

	"go.uber.org/zap"
	"golang.org/x/net/context"
	"google.golang.org/grpc"

	"github.com/imorph/gate-keeper/pkg/api/gatekeeper"
)

func TestSimpleCheck(t *testing.T) {
	logger, _ := zap.NewProduction()

	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 10, 5, 3)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			t.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)
	ok, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
	if err != nil {
		t.Error("Error when calling Check:", err)
	}
	if !ok.GetOk() {
		t.Error("Want:", true, "from check, got:", ok.GetOk())
	}

	s.Stop()

}

func TestBanByPass(t *testing.T) {
	logger, _ := zap.NewProduction()
	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 10, 5, 3)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			t.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)

	for i := 0; i < 3; i++ {
		ok, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
		if err != nil {
			t.Error("Error when calling Check:", err)
		}
		if !ok.GetOk() {
			t.Error("Want:", true, "from check, got:", ok.GetOk())
		}
	}
	ok, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
	if err == nil {
		t.Error("NO Error when calling Check, Want err=Password max attempts reached")
	}
	if ok.GetOk() {
		t.Error("Want:", false, "from check, got:", ok.GetOk())
	}

	s.Stop()

}

func TestBanByBlackList(t *testing.T) {
	logger, _ := zap.NewProduction()

	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 10, 5, 3)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			t.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)

	okB, err := c.BlackList(context.Background(), &gatekeeper.BlackListRequest{Subnet: "192.168.1.0/24", Isadd: true})
	if err != nil {
		t.Error("Error when calling BlackList:", err)
	}
	if !okB.GetOk() {
		t.Error("Want:", true, "from BlackList, got:", okB.GetOk())
	}

	okC, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "192.168.1.1"})
	if err == nil {
		t.Error("NO Error when calling Check, Want err=IP address in black-list")
	}
	if okC.GetOk() {
		t.Error("Want:", false, "from check, got:", okC.GetOk())
	}

	s.Stop()

}

func TestNoBanByPassWhiteList(t *testing.T) {
	logger, _ := zap.NewProduction()

	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 10, 5, 3)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			t.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)

	okW, err := c.WhiteList(context.Background(), &gatekeeper.WhiteListRequest{Subnet: "192.168.1.0/24", Isadd: true})
	if err != nil {
		t.Error("Error when calling WhiteList:", err)
	}
	if !okW.GetOk() {
		t.Error("Want:", true, "from WiteList, got:", okW.GetOk())
	}

	for i := 0; i < 3; i++ {
		ok, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "192.168.1.1"})
		if err != nil {
			t.Error("Error when calling Check:", err)
		}
		if !ok.GetOk() {
			t.Error("Want:", true, "from check, got:", ok.GetOk())
		}
	}
	ok, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "192.168.1.1"})
	if err != nil {
		t.Error("Error when calling Check, Want approowed")
	}
	if !ok.GetOk() {
		t.Error("Want:", true, "from check, got:", ok.GetOk())
	}

	s.Stop()

}

func TestBanByPassThenReset(t *testing.T) {
	logger, _ := zap.NewProduction()

	host := "127.0.0.1:20002"
	s := NewGateKeeperServer(host, logger, 10, 3, 5)
	go func(s *GateKeeperServer) {
		err := s.Start()
		if err != nil {
			t.Error("cant start server")
		}
	}(s)

	time.Sleep(time.Second)

	var conn *grpc.ClientConn
	conn, err := grpc.Dial(host, grpc.WithInsecure())
	if err != nil {
		t.Error("did not connect:", err)
	}
	defer func(conn *grpc.ClientConn) {
		if err := conn.Close(); err != nil {
			t.Error("Error when calling CLose:", err)
		}

	}(conn)
	c := gatekeeper.NewGateKeeperClient(conn)

	for i := 0; i < 3; i++ {
		ok, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
		if err != nil {
			t.Error("Error when calling Check:", err)
		}
		if !ok.GetOk() {
			t.Error("Want:", true, "from check, got:", ok.GetOk())
		}
	}
	ok, err := c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
	if err == nil {
		t.Error("NO Error when calling Check, Want err=Password max attempts reached")
	}
	if ok.GetOk() {
		t.Error("Want:", false, "from check, got:", ok.GetOk())
	}

	okR, err := c.Reset(context.Background(), &gatekeeper.ResetRequest{Login: "test", Ip: "127.0.0.1"})
	if err != nil {
		t.Error("Error when calling Reset:", err)
	}
	if !okR.GetOk() {
		t.Error("Want:", true, "from Reset, got:", okR.GetOk())
	}

	ok, err = c.Check(context.Background(), &gatekeeper.CheckRequest{Login: "test", Password: "test", Ip: "127.0.0.1"})
	if err != nil {
		t.Error("Error when calling Check, Want OK")
	}
	if !ok.GetOk() {
		t.Error("Want:", true, "from check, got:", ok.GetOk())
	}

	s.Stop()

}
