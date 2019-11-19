package server

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	pb "github.com/imorph/gate-keeper/pkg/api/gatekeeper"
)

type GateKeeperServer struct {
	listenHost string
	logger     *zap.Logger
}

func NewGateKeeperServer(listenHost string, logger *zap.Logger) *GateKeeperServer {
	return &GateKeeperServer{
		listenHost: listenHost,
		logger:     logger,
	}
}

func (s *GateKeeperServer) Start() error {
	lis, err := net.Listen("tcp", s.listenHost)
	if err != nil {
		s.logger.Fatal("Can not listen on", zap.String("host:port", s.listenHost), zap.Error(err))
		return err
	}
	var opts []grpc.ServerOption
	// if *tls {
	// 	if *certFile == "" {
	// 		*certFile = testdata.Path("server1.pem")
	// 	}
	// 	if *keyFile == "" {
	// 		*keyFile = testdata.Path("server1.key")
	// 	}
	// 	creds, err := credentials.NewServerTLSFromFile(*certFile, *keyFile)
	// 	if err != nil {
	// 		log.Fatalf("Failed to generate credentials %v", err)
	// 	}
	// 	opts = []grpc.ServerOption{grpc.Creds(creds)}
	// }
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterGateKeeperServer(grpcServer, s)
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *GateKeeperServer) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckReply, error) {
	s.logger.Warn("Method Check called for", zap.String("IP:", req.Ip), zap.String("Login", req.Login))
	var rep *pb.CheckReply
	ipTMP := net.ParseIP(req.Ip)
	if ipTMP == nil {
		rep.Ok = false
		s.logger.Warn("Method Check ", zap.String("This is not walid IP:", req.Ip))
		return rep, status.Errorf(codes.InvalidArgument, "IP address is malformed")
	}
	rep.Ok = true
	return rep, nil
}

func (s *GateKeeperServer) Reset(ctx context.Context, req *pb.ResetRequest) (*pb.ResetReply, error) {
	var rep *pb.ResetReply
	ipTMP := net.ParseIP(req.Ip)
	if ipTMP == nil {
		rep.Ok = false
		s.logger.Warn("Method Check ", zap.String("This is not walid IP:", req.Ip))
		return rep, status.Errorf(codes.InvalidArgument, "IP address is malformed")
	}
	rep.Ok = true
	return rep, nil
}
func (s *GateKeeperServer) WhiteList(ctx context.Context, req *pb.WhiteListRequest) (*pb.WhiteListReply, error) {
	var rep *pb.WhiteListReply
	_, _, err := net.ParseCIDR(req.Subnet)
	if err != nil {
		rep.Ok = false
		s.logger.Warn("Method Check ", zap.String("This is not walid Subnet:", req.Subnet))
		return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
	}
	rep.Ok = true
	return rep, nil
}
func (s *GateKeeperServer) BlackList(ctx context.Context, req *pb.BlackListRequest) (*pb.BlackListReply, error) {
	var rep *pb.BlackListReply
	_, _, err := net.ParseCIDR(req.Subnet)
	if err != nil {
		rep.Ok = false
		s.logger.Warn("Method Check ", zap.String("This is not walid Subnet:", req.Subnet))
		return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
	}
	rep.Ok = true
	return rep, nil
}
