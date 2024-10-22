package server

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"

	config "github.com/NStegura/saga/orders/config/server"
	"github.com/NStegura/saga/orders/internal/errs"
	"github.com/NStegura/saga/orders/internal/services/order/models"
	mock_server "github.com/NStegura/saga/orders/mocks/app/server"
	"github.com/NStegura/saga/orders/monitoring/logger"
	pb "github.com/NStegura/saga/orders/pkg/api"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

type testHelper struct {
	ctrl       *gomock.Controller
	conn       *grpc.ClientConn
	server     *GRPCServer
	mockOrder  *mock_server.MockOrder
	mockSystem *mock_server.MockSystem
}

func initTestHelper(t *testing.T) *testHelper {
	t.Helper()
	lis = bufconn.Listen(bufSize)
	ctrl := gomock.NewController(t)
	log, err := logger.Init("info")
	assert.NoError(t, err)
	mockSystem := mock_server.NewMockSystem(ctrl)
	mockPayment := mock_server.NewMockOrder(ctrl)
	cfg := config.Server{
		GRPCAddr: ":0",
	}

	server, err := New(cfg, mockPayment, mockSystem, log)
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
		ctrl:       ctrl,
		conn:       conn,
		server:     server,
		mockOrder:  mockPayment,
		mockSystem: mockSystem,
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

func TestGRPCServer_CreateOrder_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var userID int64 = 1
	var orderID int64 = 1
	description := "test order"
	orderProducts := []*pb.OrderProduct{
		{ProductId: 1, Count: 2},
	}

	th.mockOrder.EXPECT().CreateOrder(
		gomock.Any(), userID, description, gomock.Any()).Return(orderID, nil)

	client := pb.NewOrdersApiClient(th.conn)
	resp, err := client.CreateOrder(context.Background(), &pb.OrderIn{
		UserId:        userID,
		Description:   description,
		OrderProducts: orderProducts,
	})

	assert.NoError(t, err)
	assert.Equal(t, orderID, resp.OrderId)
}

func TestGRPCServer_CreateOrder_NullData(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var userID int64 = 1
	description := "test order"

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.CreateOrder(context.Background(), &pb.OrderIn{
		UserId:        userID,
		Description:   description,
		OrderProducts: nil,
	})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "null data")
}

func TestGRPCServer_CreateOrder_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var userID int64 = 1
	description := "test order"
	orderProducts := []*pb.OrderProduct{
		{ProductId: 1, Count: 2},
	}

	th.mockOrder.EXPECT().CreateOrder(
		gomock.Any(), userID, description, gomock.Any()).Return(int64(0), errors.New("internal error"))

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.CreateOrder(context.Background(), &pb.OrderIn{
		UserId:        userID,
		Description:   description,
		OrderProducts: orderProducts,
	})
	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetOrder_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var orderID int64 = 1
	state := models.ORDER_CREATED
	description := "test order"
	orderProducts := []models.OrderProduct{
		{ProductID: 1, Count: 2},
	}

	th.mockOrder.EXPECT().GetOrder(gomock.Any(), orderID).Return(models.Order{
		OrderInfo: models.OrderInfo{
			OrderID:     orderID,
			Description: description,
			State:       state,
		},
		OrderProducts: orderProducts,
	}, nil)

	client := pb.NewOrdersApiClient(th.conn)
	resp, err := client.GetOrder(context.Background(), &pb.OrderId{OrderId: orderID})

	assert.NoError(t, err)
	assert.Equal(t, orderID, resp.OrderId)
	assert.Equal(t, description, resp.Description)
	assert.Equal(t, string(state), resp.State)
}

func TestGRPCServer_GetOrder_NotFound(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var orderID int64 = 1

	th.mockOrder.EXPECT().GetOrder(gomock.Any(), orderID).Return(models.Order{}, errs.ErrNotFound)

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.GetOrder(context.Background(), &pb.OrderId{OrderId: orderID})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestGRPCServer_GetOrder_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var orderID int64 = 1

	th.mockOrder.EXPECT().GetOrder(gomock.Any(), orderID).Return(models.Order{}, errors.New("internal"))

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.GetOrder(context.Background(), &pb.OrderId{OrderId: orderID})

	assert.Error(t, err)

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetOrderStates_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var orderID int64 = 1
	orderStates := []models.State{
		{
			State:     models.ORDER_CREATED,
			CreatedAt: time.Now(),
		},
		{
			State:     models.RESERVE_CREATED,
			CreatedAt: time.Now().Add(10 * time.Minute),
		},
	}
	th.mockOrder.EXPECT().GetOrderStates(gomock.Any(), orderID).Return(orderStates, nil)

	client := pb.NewOrdersApiClient(th.conn)
	resp, err := client.GetOrderStates(context.Background(), &pb.OrderId{OrderId: orderID})

	assert.NoError(t, err)
	assert.Len(t, resp.OrderStates, 2)
	assert.Equal(t, string(models.ORDER_CREATED), resp.OrderStates[0].State)
	assert.Equal(t, string(models.RESERVE_CREATED), resp.OrderStates[1].State)
}

func TestGRPCServer_GetOrderStates_NotFound(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var orderID int64 = 1

	th.mockOrder.EXPECT().GetOrderStates(gomock.Any(), orderID).Return(nil, errs.ErrNotFound)

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.GetOrderStates(context.Background(), &pb.OrderId{OrderId: orderID})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGRPCServer_GetOrderStates_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var orderID int64 = 1

	th.mockOrder.EXPECT().GetOrderStates(gomock.Any(), orderID).Return(nil, errors.New("internal error"))

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.GetOrderStates(context.Background(), &pb.OrderId{OrderId: orderID})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "internal error")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetOrders_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var userID int64 = 1
	orders := []models.OrderInfo{
		{
			OrderID:     1,
			Description: "Test order 1",
			State:       models.ORDER_CREATED,
		},
		{
			OrderID:     2,
			Description: "Test order 2",
			State:       models.ORDER_CREATED,
		},
	}
	th.mockOrder.EXPECT().GetOrders(gomock.Any(), userID).Return(orders, nil)

	// Выполняем запрос
	client := pb.NewOrdersApiClient(th.conn)
	resp, err := client.GetOrders(context.Background(), &pb.UserId{UserId: userID})

	assert.NoError(t, err)
	assert.Len(t, resp.Orders, 2)
	assert.Equal(t, orders[0].OrderID, resp.Orders[0].OrderId)
	assert.Equal(t, orders[0].Description, resp.Orders[0].Description)
	assert.Equal(t, string(orders[0].State), resp.Orders[0].State)
	assert.Equal(t, orders[1].OrderID, resp.Orders[1].OrderId)
	assert.Equal(t, orders[1].Description, resp.Orders[1].Description)
	assert.Equal(t, string(orders[1].State), resp.Orders[1].State)
}

func TestGRPCServer_GetOrders_NotFound(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var userID int64 = 1
	th.mockOrder.EXPECT().GetOrders(gomock.Any(), userID).Return(nil, errs.ErrNotFound)

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.GetOrders(context.Background(), &pb.UserId{UserId: userID})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
}

func TestGRPCServer_GetOrders_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	var userID int64 = 1
	th.mockOrder.EXPECT().GetOrders(gomock.Any(), userID).Return(nil, errors.New("internal"))

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.GetOrders(context.Background(), &pb.UserId{UserId: userID})

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "internal")

	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.Internal, st.Code())
}

func TestGRPCServer_GetPing_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockSystem.EXPECT().Ping(gomock.Any()).Return(nil)

	client := pb.NewOrdersApiClient(th.conn)
	res, err := client.GetPing(context.Background(), &emptypb.Empty{})
	assert.NoError(t, err)
	assert.True(t, res.Pong)
}

func TestGRPCServer_GetPing_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockSystem.EXPECT().Ping(gomock.Any()).Return(errors.New("internal error"))

	client := pb.NewOrdersApiClient(th.conn)
	_, err := client.GetPing(context.Background(), &emptypb.Empty{})
	assert.Error(t, err)
}
