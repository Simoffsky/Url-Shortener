package services

import (
	"context"
	"url-shorter/internal/kafka"
	"url-shorter/internal/models"
	pb "url-shorter/pkg/proto/stats"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type StatService interface {
	GetStatForLink(short string) (*models.LinkStat, error)
	SendStat(stat *models.LinkStatVisitor) error
}

type StatServiceGRPC struct {
	statsClient            pb.StatServiceClient
	messageProducerManager *kafka.ProducerManager
}

func NewStatServiceGRPC(statsAddr string, messageProducerManager *kafka.ProducerManager) (*StatServiceGRPC, error) {
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	conn, err := grpc.NewClient(statsAddr, opts...)
	if err != nil {
		return nil, err
	}

	client := pb.NewStatServiceClient(conn)

	return &StatServiceGRPC{
		statsClient:            client,
		messageProducerManager: messageProducerManager,
	}, nil
}

func (s *StatServiceGRPC) GetStatForLink(short string) (*models.LinkStat, error) {
	resp, err := s.statsClient.GetStatForLink(context.Background(), &pb.LinkStatRequest{Short: short})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return nil, models.ErrLinkNotFound
		}
		return nil, err
	}

	return &models.LinkStat{
		ShortLink:      resp.ShortLink,
		Clicks:         int(resp.Clicks),
		LastAccessedAt: resp.LastAccessedAt,
		UniqueClicks:   int(resp.UniqueClicks),
	}, nil
}

func (s *StatServiceGRPC) SendStat(stat *models.LinkStatVisitor) error {
	go func() {
		s.messageProducerManager.C() <- *stat
	}()
	return nil
}
