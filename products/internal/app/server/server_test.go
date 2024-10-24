package server

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"

	config "github.com/NStegura/saga/products/config/server"
	"github.com/NStegura/saga/products/internal/errs"
	productModels "github.com/NStegura/saga/products/internal/services/product/models"
	mock_server "github.com/NStegura/saga/products/mocks/app/server"
	"github.com/NStegura/saga/products/monitoring/logger"
	pb "github.com/NStegura/saga/products/pkg/api"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type testHelper struct {
	ctrl        *gomock.Controller
	conn        *grpc.ClientConn
	server      *GRPCServer
	mockProduct *mock_server.MockProduct
	mockSystem  *mock_server.MockSystem
}

func initTestHelper(t *testing.T) *testHelper {
	t.Helper()
	lis = bufconn.Listen(bufSize)
	ctrl := gomock.NewController(t)
	log, _ := logger.Init("info")
	mockSystem := mock_server.NewMockSystem(ctrl)
	mockProduct := mock_server.NewMockProduct(ctrl)
	cfg := config.Server{
		GRPCAddr: ":0",
	}

	server, err := New(cfg, mockProduct, mockSystem, log)
	assert.NoError(t, err)

	go func() {
		if err := server.server.Serve(lis); err != nil {
			t.Errorf("Failed to serve: %v", err)
			return
		}
	}()

	conn, err := grpc.DialContext(
		context.Background(),
		"bufconn",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	assert.NoError(t, err)

	if conn.GetState() != connectivity.Ready {
		conn.Connect()
	}

	return &testHelper{
		ctrl:        ctrl,
		conn:        conn,
		server:      server,
		mockProduct: mockProduct,
		mockSystem:  mockSystem,
	}
}

func (th *testHelper) finish() {
	_ = th.conn.Close()
	_ = th.server.Shutdown(context.Background())
	th.ctrl.Finish()
}

func TestGRPCServer_GetName(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()
	assert.Equal(t, "grpc server", th.server.Name())
}

func TestGRPCServer_GetProducts__Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockProduct.EXPECT().GetProducts(gomock.Any()).Return([]productModels.Product{
		{
			ProductID:   1,
			Category:    "Electronics",
			Name:        "Laptop",
			Description: "A powerful laptop",
			Count:       10,
		},
	}, nil)

	client := pb.NewProductsApiClient(th.conn)

	// Тестирование GetProducts
	res, err := client.GetProducts(context.Background(), &empty.Empty{})
	assert.NoError(t, err)
	assert.Len(t, res.Products, 1)
	assert.Equal(t, int64(1), res.Products[0].ProductId)
	assert.Equal(t, "Laptop", res.Products[0].Name)
	assert.Equal(t, "Electronics", res.Products[0].Category)
}

func TestGRPCServer_GetProducts__NotFound(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockProduct.EXPECT().GetProducts(gomock.Any()).Return(nil, errs.ErrNotFound)
	client := pb.NewProductsApiClient(th.conn)

	_, err := client.GetProducts(context.Background(), &empty.Empty{})
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGRPCServer_GetProducts_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockProduct.EXPECT().GetProducts(gomock.Any()).Return(nil, errors.New("internal error"))
	client := pb.NewProductsApiClient(th.conn)

	_, err := client.GetProducts(context.Background(), &empty.Empty{})
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetProductInfo__Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockProduct.EXPECT().GetProductInfo(gomock.Any(), int64(1)).Return(productModels.Product{
		ProductID:   1,
		Category:    "Electronics",
		Name:        "Laptop",
		Description: "A powerful laptop",
		Count:       10,
	}, nil)

	client := pb.NewProductsApiClient(th.conn)

	req := &pb.ProductId{ProductId: 1}
	res, err := client.GetProductInfo(context.Background(), req)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), res.ProductId)
	assert.Equal(t, "Laptop", res.Name)
	assert.Equal(t, "Electronics", res.Category)
}

func TestGRPCServer_GetProductInfo__NotFound(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockProduct.EXPECT().GetProductInfo(gomock.Any(), int64(1)).Return(productModels.Product{}, errs.ErrNotFound)
	client := pb.NewProductsApiClient(th.conn)

	req := &pb.ProductId{ProductId: 1}
	_, err := client.GetProductInfo(context.Background(), req)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGRPCServer_GetProductInfo__InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockProduct.EXPECT().GetProductInfo(
		gomock.Any(), int64(1)).Return(productModels.Product{}, errors.New("internal error"))
	client := pb.NewProductsApiClient(th.conn)

	req := &pb.ProductId{ProductId: 1}
	_, err := client.GetProductInfo(context.Background(), req)
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetPing_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockSystem.EXPECT().Ping(gomock.Any()).Return(nil)

	client := pb.NewProductsApiClient(th.conn)

	res, err := client.GetPing(context.Background(), &empty.Empty{})
	assert.NoError(t, err)
	assert.True(t, res.Pong)
}

func TestGRPCServer_GetPing_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockSystem.EXPECT().Ping(gomock.Any()).Return(errors.New("internal error"))

	client := pb.NewProductsApiClient(th.conn)

	_, err := client.GetPing(context.Background(), &empty.Empty{})
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}
