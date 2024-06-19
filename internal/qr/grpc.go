package qr

import (
	"context"
	pb "url-shorter/pkg/proto/qr"
)

func (s *QRServer) GetQRCode(ctx context.Context, in *pb.QRCodeRequest) (*pb.QRCodeResponse, error) {
	qr, err := s.service.GetQRCode(in.Url, int(in.Size))
	if err != nil {
		return nil, err
	}
	return &pb.QRCodeResponse{QrCode: qr}, nil
}
