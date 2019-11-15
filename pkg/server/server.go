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
}

func NewGateKeeperServer() *GateKeeperServer {
	return &GateKeeperServer{}
}

func (s *GateKeeperServer) Start(listenHost string, logger *zap.Logger) error {
	lis, err := net.Listen("tcp", listenHost)
	if err != nil {
		logger.Fatal("Can not listen on", zap.String("host:port", listenHost), zap.Error(err))
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
	pb.RegisterGateKeeperServer(grpcServer, NewGateKeeperServer())
	err = grpcServer.Serve(lis)
	if err != nil {
		return err
	}
	return nil
}

func (s *GateKeeperServer) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Check not implemented")
}

func (s *GateKeeperServer) Reset(ctx context.Context, req *pb.ResetRequest) (*pb.ResetReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Reset not implemented")
}
func (s *GateKeeperServer) WhiteList(ctx context.Context, req *pb.WhiteListRequest) (*pb.WhiteListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WhiteList not implemented")
}
func (s *GateKeeperServer) BlackList(ctx context.Context, req *pb.BlackListRequest) (*pb.BlackListReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BlackList not implemented")
}
