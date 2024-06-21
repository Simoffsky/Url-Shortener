package stats

import (
	"context"
	pb "url-shorter/pkg/proto/stats"
)

func (s *StatsServer) GetStatForLink(ctx context.Context, in *pb.LinkStatRequest) (*pb.LinkStatResponse, error) {
	stat, err := s.service.GetStatForLink(in.Short)
	if err != nil {
		return nil, err
	}

	return &pb.LinkStatResponse{
		ShortLink:      stat.ShortLink,
		Clicks:         int32(stat.Clicks),
		LastAccessedAt: stat.LastAccessedAt,
		UniqueClicks:   int32(stat.UniqueClicks),
	}, nil
}
