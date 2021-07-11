package tunnels

import (
	"context"
	"log"

	"github.com/tnynlabs/wyrm/pkg/tunnels/protobuf"
	"github.com/tnynlabs/wyrm/pkg/utils"
	"google.golang.org/grpc"
)

type Service interface {
	InvokeDevice(deviceID int64, pattern string, data string) (*InvokeResponse, error)
	RevokeDevice(deviceID int64)
}

type httpGrpcService struct {
	client protobuf.TunnelManagerClient
}

//"123.0.0.01.1:9090"
func CreateHttpGrpcService(target string) Service {
	conn, err := grpc.Dial(target, grpc.WithInsecure())
	if err != nil {
		log.Printf("GRPC Server Error")
	}
	client := protobuf.NewTunnelManagerClient(conn)
	return &httpGrpcService{client}
}

func (s *httpGrpcService) InvokeDevice(deviceID int64, pattern string, data string) (*InvokeResponse, error) {
	invokeRequest := protobuf.InvokeRequest{
		DeviceId: deviceID,
		Pattern:  pattern,
		Data:     data,
	}
	invokeResp, err := s.client.InvokeDevice(context.Background(), &invokeRequest)
	if err != nil {
		return nil, &utils.ServiceErr{
			Code:    ConnectionErrorCode,
			Message: "Invalid Key",
		}
	}
	return &InvokeResponse{Data: invokeResp.Data}, nil
}

func (s *httpGrpcService) RevokeDevice(deviceID int64) {

}

type InvokeResponse struct {
	Data string
}
