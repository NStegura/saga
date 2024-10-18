package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"

	"github.com/NStegura/saga/payments/internal/services/payment/models"

	config "github.com/NStegura/saga/payments/config/server"
	pb "github.com/NStegura/saga/payments/pkg/api"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GRPCServer struct {
	pb.UnimplementedPaymentsApiServer

	server   *grpc.Server
	cfg      config.Server
	payments Payments
	system   System

	logger *logrus.Logger
}

func New(cfg config.Server, p Payments, s System, logger *logrus.Logger, opts ...grpc.ServerOption) (*GRPCServer, error) {
	grpcServer := grpc.NewServer(opts...)

	server := &GRPCServer{
		server:   grpcServer,
		cfg:      cfg,
		payments: p,
		system:   s,
		logger:   logger,
	}
	pb.RegisterPaymentsApiServer(grpcServer, server)
	return server, nil
}

func (s *GRPCServer) Start(_ context.Context) error {
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

	if err = s.server.Serve(lis); err != nil {
		return fmt.Errorf("failed to start grpc server: %w", err)
	}
	return nil
}

func (s *GRPCServer) Shutdown(ctx context.Context) (err error) {
	doneCh := make(chan struct{})

	go func() {
		s.server.GracefulStop()
		close(doneCh)
	}()

	select {
	case <-ctx.Done():
		log.Println("Shutdown timeout reached, force closing.")
		s.server.Stop()
		err = ctx.Err()
	case <-doneCh:
		log.Println("Shutdown success")
	}
	return
}

func (s *GRPCServer) Name() string {
	return "grpc server"
}

func (s *GRPCServer) UpdatePaymentStatus(ctx context.Context, req *pb.PayStatus) (*empty.Empty, error) {
	var paymentStatus models.PaymentMessageStatus
	if req.Status {
		paymentStatus = models.COMPLETED
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
