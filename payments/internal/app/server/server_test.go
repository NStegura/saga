package server

import (
	"context"
	"errors"
	"net"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/emptypb"

	config "github.com/NStegura/saga/payments/config/server"
	"github.com/NStegura/saga/payments/internal/monitoring/logger"
	"github.com/NStegura/saga/payments/internal/services/payment/models"
	mock_server "github.com/NStegura/saga/payments/mocks/app/server"
	pb "github.com/NStegura/saga/payments/pkg/api"
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
	mockPayment *mock_server.MockPayments
	mockSystem  *mock_server.MockSystem
}

func initTestHelper(t *testing.T) *testHelper {
	t.Helper()
	lis = bufconn.Listen(bufSize)
	ctrl := gomock.NewController(t)
	log, err := logger.Init("info")
	assert.NoError(t, err)
	mockSystem := mock_server.NewMockSystem(ctrl)
	mockPayment := mock_server.NewMockPayments(ctrl)
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
		ctrl:        ctrl,
		conn:        conn,
		server:      server,
		mockPayment: mockPayment,
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

func TestGRPCServer_UpdatePaymentStatus_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	tests := []struct {
		name           string
		req            *pb.PayStatus
		expectedStatus models.PaymentMessageStatus
	}{
		{
			name: "Ok, pay true",
			req: &pb.PayStatus{
				OrderId: 123,
				Status:  true,
			},
			expectedStatus: models.COMPLETED,
		},
		{
			name: "Ok, pay false",
			req: &pb.PayStatus{
				OrderId: 123,
				Status:  false,
			},
			expectedStatus: models.FAILED,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := pb.NewPaymentsApiClient(th.conn)
			th.mockPayment.EXPECT().UpdatePaymentStatus(
				gomock.Any(), test.req.OrderId, test.expectedStatus).Return(nil)
			_, err := client.UpdatePaymentStatus(context.Background(), test.req)
			assert.NoError(t, err)
		})
	}
}

func TestGRPCServer_UpdatePaymentStatus_Failed(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockPayment.EXPECT().UpdatePaymentStatus(
		gomock.Any(), gomock.Any(), gomock.Any()).Return(errors.New("internal"))

	client := pb.NewPaymentsApiClient(th.conn)
	_, err := client.UpdatePaymentStatus(context.Background(), &pb.PayStatus{
		OrderId: 123,
		Status:  false,
	})
	assert.Error(t, err)
}

func TestGRPCServer_GetPing_Success(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockSystem.EXPECT().Ping(gomock.Any()).Return(nil)

	client := pb.NewPaymentsApiClient(th.conn)
	res, err := client.GetPing(context.Background(), &emptypb.Empty{})
	assert.NoError(t, err)
	assert.True(t, res.Pong)
}

func TestGRPCServer_GetPing_InternalError(t *testing.T) {
	th := initTestHelper(t)
	defer th.finish()

	th.mockSystem.EXPECT().Ping(gomock.Any()).Return(errors.New("internal error"))

	client := pb.NewPaymentsApiClient(th.conn)
	_, err := client.GetPing(context.Background(), &emptypb.Empty{})
	assert.Error(t, err)
}
