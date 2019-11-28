package server

import (
	"context"
	"net"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"

	pb "github.com/imorph/gate-keeper/pkg/api/gatekeeper"
	"github.com/imorph/gate-keeper/pkg/core"
)

// GateKeeperServer is basic struct for server
type GateKeeperServer struct {
	listenHost string
	logger     *zap.Logger
	ip         *core.LimitersCache
	login      *core.LimitersCache
	pass       *core.LimitersCache
	white      *core.List
	black      *core.List
	grpcServ   *grpc.Server
}

// NewGateKeeperServer accept settings and returns new gatekeeper server
func NewGateKeeperServer(listenHost string, logger *zap.Logger, ipMax, loginMax, passMax int) *GateKeeperServer {
	var opts []grpc.ServerOption
	return &GateKeeperServer{
		listenHost: listenHost,
		logger:     logger,
		ip:         core.NewCache(ipMax, time.Minute, 2*time.Minute),
		login:      core.NewCache(loginMax, time.Minute, 2*time.Minute),
		pass:       core.NewCache(passMax, time.Minute, 2*time.Minute),
		white:      core.NewList(),
		black:      core.NewList(),
		grpcServ:   grpc.NewServer(opts...),
	}
}

// Start starts new server
func (s *GateKeeperServer) Start() error {
	lis, err := net.Listen("tcp", s.listenHost)
	if err != nil {
		s.logger.Fatal("Can not listen on", zap.String("host:port", s.listenHost), zap.Error(err))
		return err
	}

	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			s.ip.HouseKeep()
			s.logger.Debug("Housekeeping for IPs done.")
		}
	}()
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			s.login.HouseKeep()
			s.logger.Debug("Housekeeping for Logins done.")
		}
	}()
	go func() {
		ticker := time.NewTicker(1 * time.Minute)
		for range ticker.C {
			s.pass.HouseKeep()
			s.logger.Debug("Housekeeping for Passwords done.")
		}
	}()

	pb.RegisterGateKeeperServer(s.grpcServ, s)
	err = s.grpcServ.Serve(lis)
	if err != nil {
		return err
	}

	return nil
}

// Stop tries to gracefully stop server
func (s *GateKeeperServer) Stop() {
	s.logger.Debug("stopping server")
	s.grpcServ.GracefulStop()
}

// Check is handler for external request to check ip/login/pass bruteforcing
func (s *GateKeeperServer) Check(ctx context.Context, req *pb.CheckRequest) (*pb.CheckReply, error) {
	s.logger.Debug("Method Check called for", zap.String("IP:", req.Ip), zap.String("Login", req.Login))
	var rep *pb.CheckReply
	ipTMP := net.ParseIP(req.Ip)
	if ipTMP == nil {
		rep = &pb.CheckReply{
			Ok: false,
		}
		s.logger.Warn("Method Check ", zap.String("This is not valid IP:", req.Ip))
		return rep, status.Errorf(codes.InvalidArgument, "IP address is malformed")
	}
	ok, err := s.black.LookUpIP(req.GetIp())
	if err != nil {
		rep = &pb.CheckReply{
			Ok: false,
		}
		s.logger.Warn("Method Check ", zap.String("This is not valid IP:", req.Ip))
		return rep, status.Errorf(codes.InvalidArgument, "IP address is malformed")
	}
	if ok {
		rep = &pb.CheckReply{
			Ok: false,
		}
		s.logger.Debug("Method Check ", zap.String("This IP in BLACK LIST:", req.Ip))
		return rep, status.Errorf(codes.PermissionDenied, "IP address in black-list")
	}
	ok, err = s.white.LookUpIP(req.GetIp())
	if err != nil {
		rep = &pb.CheckReply{
			Ok: false,
		}
		s.logger.Warn("Method Check ", zap.String("This is not valid IP:", req.Ip))
		return rep, status.Errorf(codes.InvalidArgument, "IP address is malformed")
	}
	if ok {
		rep = &pb.CheckReply{
			Ok: true,
		}
		s.logger.Debug("Method Check ", zap.String("This IP in WHITE LIST:", req.Ip))
		return rep, nil
	}

	okIP := s.ip.Check(req.GetIp())
	okLogin := s.login.Check(req.GetLogin())
	okPass := s.pass.Check(req.GetPassword())
	if !okIP {
		rep = &pb.CheckReply{
			Ok: false,
		}
		s.logger.Debug("Method Check ", zap.String("Too many attempts for IP", req.Ip))
		return rep, status.Errorf(codes.PermissionDenied, "IP address max attempts reached")
	}
	if !okLogin {
		rep = &pb.CheckReply{
			Ok: false,
		}
		s.logger.Debug("Method Check ", zap.String("Too many attempts for Login:", req.Login))
		return rep, status.Errorf(codes.PermissionDenied, "Login max attempts reached")
	}
	if !okPass {
		rep = &pb.CheckReply{
			Ok: false,
		}
		s.logger.Debug("Method Check ", zap.String("Too many attempts for Pass:", "***"))
		return rep, status.Errorf(codes.PermissionDenied, "Password max attempts reached")
	}
	s.logger.Debug("Method Check successfully executed for", zap.String("IP:", req.Ip), zap.String("Login", req.Login))
	rep = &pb.CheckReply{
		Ok: true,
	}
	return rep, nil
}

// Reset is handler for reseting counters for particular IP/Login
func (s *GateKeeperServer) Reset(ctx context.Context, req *pb.ResetRequest) (*pb.ResetReply, error) {
	s.logger.Debug("Method Reset called for", zap.String("IP:", req.Ip), zap.String("Login", req.Login))
	var rep *pb.ResetReply
	ipTMP := net.ParseIP(req.Ip)
	if ipTMP == nil {
		rep = &pb.ResetReply{
			Ok: false,
		}
		s.logger.Warn("Method Check ", zap.String("This is not valid IP:", req.Ip))
		return rep, status.Errorf(codes.InvalidArgument, "IP address is malformed")
	}
	s.ip.Reset(req.GetIp())
	s.login.Reset(req.GetLogin())
	s.logger.Debug("Method Reset successfully executed for", zap.String("IP:", req.Ip), zap.String("Login", req.Login))
	rep = &pb.ResetReply{
		Ok: true,
	}
	return rep, nil
}

// WhiteList is handler for adding CIDR to white list
func (s *GateKeeperServer) WhiteList(ctx context.Context, req *pb.WhiteListRequest) (*pb.WhiteListReply, error) {
	s.logger.Debug("Method WhiteList called for", zap.String("IP:", req.Subnet), zap.Bool("Add to list", req.Isadd))
	var rep *pb.WhiteListReply
	_, _, err := net.ParseCIDR(req.Subnet)
	if err != nil {
		rep = &pb.WhiteListReply{
			Ok: false,
		}
		s.logger.Warn("Method Check ", zap.String("This is not valid Subnet:", req.Subnet))
		return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
	}
	if req.GetIsadd() {
		err := s.white.InsertCIDR(req.GetSubnet())
		if err != nil {
			rep = &pb.WhiteListReply{
				Ok: false,
			}
			s.logger.Warn("Method Check ", zap.String("This is not valid Subnet:", req.Subnet))
			return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
		}
	} else {
		err := s.white.DeleteCIDR(req.GetSubnet())
		if err != nil {
			rep = &pb.WhiteListReply{
				Ok: false,
			}
			s.logger.Warn("Method Check ", zap.String("This is not valid Subnet:", req.Subnet))
			return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
		}
	}
	s.logger.Debug("Method WhiteList successfully executed for", zap.String("IP:", req.Subnet), zap.Bool("Add to list", req.Isadd))
	rep = &pb.WhiteListReply{
		Ok: true,
	}
	return rep, nil
}

// BlackList is handler for adding CIDR to black list
func (s *GateKeeperServer) BlackList(ctx context.Context, req *pb.BlackListRequest) (*pb.BlackListReply, error) {
	s.logger.Debug("Method BlackList called for", zap.String("IP:", req.Subnet), zap.Bool("Add to list", req.Isadd))
	var rep *pb.BlackListReply
	_, _, err := net.ParseCIDR(req.Subnet)
	if err != nil {
		rep = &pb.BlackListReply{
			Ok: false,
		}
		s.logger.Warn("Method Check ", zap.String("This is not valid Subnet:", req.Subnet))
		return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
	}

	if req.GetIsadd() {
		err := s.black.InsertCIDR(req.GetSubnet())
		if err != nil {
			rep = &pb.BlackListReply{
				Ok: false,
			}
			s.logger.Warn("Method Check ", zap.String("This is not valid Subnet:", req.Subnet))
			return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
		}
	} else {
		err := s.black.DeleteCIDR(req.GetSubnet())
		if err != nil {
			rep = &pb.BlackListReply{
				Ok: false,
			}
			s.logger.Warn("Method Check ", zap.String("This is not valid Subnet:", req.Subnet))
			return rep, status.Errorf(codes.InvalidArgument, "Subnet CIDR is malformed")
		}
	}
	s.logger.Debug("Method BlackList successfully executed for", zap.String("IP:", req.Subnet), zap.Bool("Add to list", req.Isadd))
	rep = &pb.BlackListReply{
		Ok: true,
	}
	return rep, nil
}
