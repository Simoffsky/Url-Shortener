package repository

import (
	"context"
	pb "url-shorter/pkg/proto/qr"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type QrGRPCRepository struct {
	qrClient pb.QRCodeServiceClient
}

func NewQrGRPCRepository(qrAddr string) (*QrGRPCRepository, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(qrAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewQRCodeServiceClient(conn)
	return &QrGRPCRepository{
		qrClient: client,
	}, nil
}

func (r *QrGRPCRepository) GetQRCode(link string, imgSize int) ([]byte, error) {
	resp, err := r.qrClient.GetQRCode(context.Background(), &pb.QRCodeRequest{
		Url:  link,
		Size: int32(imgSize),
	})
	if err != nil {
		return nil, err
	}
	return resp.QrCode, nil
}
