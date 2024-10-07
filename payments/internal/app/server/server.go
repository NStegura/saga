package grpcserver

import (
	"context"
	"fmt"
	"github.com/NStegura/saga/payments/internal/services/payment/models"
	"github.com/golang/protobuf/ptypes/empty"
	"net"

	config "github.com/NStegura/saga/payments/config/server"
	pb "github.com/NStegura/saga/payments/pkg/api"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	pb.UnimplementedPaymentsApiServer

	cfg      *config.Server
	payments Payments
	system   System

	logger *logrus.Logger
}

func New(cfg *config.Server, p Payments, s System, logger *logrus.Logger) (*GRPCServer, error) {
	return &GRPCServer{
		cfg:      cfg,
		payments: p,
		system:   s,
		logger:   logger,
	}, nil
}

func (s *GRPCServer) Start(opts ...grpc.ServerOption) error {
	s.logger.Infof("starting GRPCServer %s", s.cfg.GRPCAddr)
	lis, err := net.Listen("tcp", s.cfg.GRPCAddr)
	defer func() {
		if err = lis.Close(); err != nil {
			s.logger.Error("failed to close listener")
		}
	}()
	if err != nil {
		return fmt.Errorf("failed to create network listener: %w", err)
	}
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterPaymentsApiServer(grpcServer, s)
	if err = grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}
	return nil
}

func (s *GRPCServer) UpdatePaymentStatus(ctx context.Context, req *pb.PayStatus) (*empty.Empty, error) {
	var paymentStatus models.PaymentMessageStatus
	if req.Status {
		paymentStatus = models.CREATED
	} else if !req.Status {
		paymentStatus = models.FAILED
	}
	if err := s.payments.UpdatePaymentStatus(ctx, req.OrderId, paymentStatus); err != nil {
		return &emptypb.Empty{}, err
	}
	return &emptypb.Empty{}, nil
}

func (s *GRPCServer) GetPing(ctx context.Context, _ *emptypb.Empty) (*pb.Pong, error) {
	if err := s.system.Ping(ctx); err != nil {
		return nil, err
	}
	return &pb.Pong{
		Pong: true,
	}, nil
}
