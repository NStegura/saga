package server

import (
	"context"
	"errors"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/NStegura/saga/orders/internal/errs"
	"github.com/NStegura/saga/orders/internal/services/order/models"

	config "github.com/NStegura/saga/orders/config/server"
	pb "github.com/NStegura/saga/orders/pkg/api"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedOrdersApiServer

	server *grpc.Server
	cfg    config.Server
	order  Order
	system System

	logger *logrus.Logger
}

func New(cfg config.Server, o Order, s System, logger *logrus.Logger, opts ...grpc.ServerOption) (*GRPCServer, error) {
	grpcServer := grpc.NewServer(opts...)

	server := &GRPCServer{
		server: grpcServer,
		cfg:    cfg,
		order:  o,
		system: s,
		logger: logger,
	}
	pb.RegisterOrdersApiServer(grpcServer, server)
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

func (s *GRPCServer) CreateOrder(ctx context.Context, req *pb.OrderIn) (*pb.OrderId, error) {
	if req.OrderProducts == nil {
		return nil, errs.ErrNullData
	}
	orderProductIn := make([]models.OrderProduct, 0, len(req.OrderProducts))
	for _, orderProduct := range req.OrderProducts {
		orderProductIn = append(orderProductIn, models.OrderProduct{
			ProductID: orderProduct.ProductId,
			Count:     orderProduct.Count,
		})
	}

	orderID, err := s.order.CreateOrder(ctx, req.UserId, req.Description, orderProductIn)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.OrderId{
		OrderId: orderID,
	}, nil
}

func (s *GRPCServer) GetOrder(ctx context.Context, req *pb.OrderId) (*pb.OrderOut, error) {
	order, err := s.order.GetOrder(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	orderProductsOut := make([]*pb.OrderProduct, 0, len(order.OrderProducts))
	for _, orderProduct := range order.OrderProducts {
		orderProductsOut = append(
			orderProductsOut,
			&pb.OrderProduct{
				ProductId: orderProduct.ProductID,
				Count:     orderProduct.Count,
			},
		)
	}

	return &pb.OrderOut{
		OrderId:       order.OrderID,
		OrderProducts: orderProductsOut,
		Description:   order.Description,
		State:         string(order.State),
	}, err
}

func (s *GRPCServer) GetOrderStates(ctx context.Context, req *pb.OrderId) (*pb.States, error) {
	orderStates, err := s.order.GetOrderStates(ctx, req.OrderId)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}
	orderStatesOut := make([]*pb.OrderState, 0, len(orderStates))
	for _, oState := range orderStates {
		orderStatesOut = append(orderStatesOut, &pb.OrderState{
			State: string(oState.State),
			Time:  timestamppb.New(oState.CreatedAt),
		})
	}
	return &pb.States{OrderStates: orderStatesOut}, nil
}

func (s *GRPCServer) GetOrders(ctx context.Context, req *pb.UserId) (*pb.Orders, error) {
	orders, err := s.order.GetOrders(ctx, req.UserId)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	outOrders := make([]*pb.OrderInfoOut, 0, len(orders))
	for _, order := range orders {
		outOrders = append(
			outOrders, &pb.OrderInfoOut{
				OrderId:     order.OrderID,
				Description: order.Description,
				State:       string(order.State),
			},
		)
	}
	return &pb.Orders{Orders: outOrders}, nil
}

func (s *GRPCServer) GetPing(ctx context.Context, _ *empty.Empty) (*pb.Pong, error) {
	if err := s.system.Ping(ctx); err != nil {
		return nil, err
	}
	return &pb.Pong{
		Pong: true,
	}, nil
}
