package server

import (
	"context"
	"errors"
	"fmt"
	"github.com/NStegura/saga/products/internal/errs"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
	"net"

	config "github.com/NStegura/saga/products/config/server"
	pb "github.com/NStegura/saga/products/pkg/api"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
)

type GRPCServer struct {
	pb.UnimplementedProductsApiServer

	server   *grpc.Server
	cfg      config.Server
	products Product
	system   System

	logger *logrus.Logger
}

func New(cfg config.Server, p Product, s System, logger *logrus.Logger, opts ...grpc.ServerOption) (*GRPCServer, error) {
	grpcServer := grpc.NewServer(opts...)

	server := &GRPCServer{
		server:   grpcServer,
		cfg:      cfg,
		products: p,
		system:   s,
		logger:   logger,
	}
	pb.RegisterProductsApiServer(grpcServer, server)
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
	return nil
}

func (s *GRPCServer) Name() string {
	return "grpc server"
}

func (s *GRPCServer) GetProducts(ctx context.Context, _ *empty.Empty) (*pb.Products, error) {
	products, err := s.products.GetProducts(ctx)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	outProducts := make([]*pb.Product, 0, len(products))
	for _, product := range products {
		outProducts = append(
			outProducts, &pb.Product{
				ProductId:   product.ProductID,
				Category:    product.Category,
				Name:        product.Name,
				Description: product.Description,
				Count:       product.Count,
			},
		)
	}
	return &pb.Products{Products: outProducts}, nil
}

func (s *GRPCServer) GetProductInfo(ctx context.Context, req *pb.ProductId) (*pb.Product, error) {
	pInfo, err := s.products.GetProductInfo(ctx, req.ProductId)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &pb.Product{
		ProductId:   pInfo.ProductID,
		Category:    pInfo.Category,
		Name:        pInfo.Name,
		Description: pInfo.Description,
		Count:       pInfo.Count,
	}, nil
}

func (s *GRPCServer) GetPing(ctx context.Context, _ *empty.Empty) (*pb.Pong, error) {
	if err := s.system.Ping(ctx); err != nil {
		return nil, err
	}
	return &pb.Pong{
		Pong: true,
	}, nil
}
