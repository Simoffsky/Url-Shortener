package stats

import (
	"context"
	"errors"
	"url-shorter/internal/models"
	pb "url-shorter/pkg/proto/stats"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *StatsServer) GetStatForLink(ctx context.Context, in *pb.LinkStatRequest) (*pb.LinkStatResponse, error) {
	stat, err := s.service.GetStatForLink(in.Short)
	if err != nil {

		if errors.Is(err, models.ErrLinkNotFound) {
			return nil, status.Error(codes.NotFound, "link not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &pb.LinkStatResponse{
		ShortLink:      stat.ShortLink,
		Clicks:         int32(stat.Clicks),
		LastAccessedAt: stat.LastAccessedAt,
		UniqueClicks:   int32(stat.UniqueClicks),
	}, nil
}
