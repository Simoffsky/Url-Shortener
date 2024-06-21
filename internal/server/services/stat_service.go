package services

import (
	"url-shorter/internal/kafka"
	"url-shorter/internal/models"
	pb "url-shorter/pkg/proto/stats"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
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
	panic("implement me (GetStatForLink)")
}

func (s *StatServiceGRPC) SendStat(stat *models.LinkStatVisitor) error {
	go func() {
		s.messageProducerManager.C() <- *stat
	}()
	return nil
}
