package qr

import (
	"context"
	"errors"
	"url-shorter/internal/models"
	pb "url-shorter/pkg/proto/qr"
)

func (s *QRServer) GetQRCode(ctx context.Context, in *pb.QRCodeRequest) (*pb.QRCodeResponse, error) {
	qr, err := s.service.GetQRCode(in.Url)
	if errors.Is(err, models.ErrLinkNotFound) {
		qr, err = s.service.CreateQRCode(in.Url, 256)
		if err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	return &pb.QRCodeResponse{QrCode: qr}, nil
}
