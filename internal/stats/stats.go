package stats

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"url-shorter/internal/config"
	"url-shorter/internal/kafka"
	"url-shorter/internal/repository/stats"
	"url-shorter/pkg/log"
	pb "url-shorter/pkg/proto/stats"

	"google.golang.org/grpc"
)

type StatsServer struct {
	pb.UnimplementedStatServiceServer
	kafkaReader *KafkaReader
	service     StatsService
	config      config.Config
	logger      log.Logger
}

func NewStatsServer(config config.Config) *StatsServer {
	return &StatsServer{
		config: config,
		logger: log.NewDefaultLogger(log.LevelFromString(config.LoggerLevel)),
	}
}

func (s *StatsServer) Start() error {
	if err := s.configure(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", s.config.StatsAddr)
	if err != nil {
		return err
	}

	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterStatServiceServer(grpcServer, s)

	ctx, cancel := context.WithCancel(context.Background())
	s.kafkaReader.Start(ctx)

	//Graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		s.logger.Info("Starting gRPC server on address: " + s.config.StatsAddr)
		if err := grpcServer.Serve(lis); err != nil {
			s.logger.Error("failed to serve gRPC: %s" + err.Error())
		}
	}()

	<-sigChan

	s.logger.Info("Shutting down gRPC server gracefully")
	stopCh := make(chan struct{})
	grpcServer.GracefulStop()
	close(stopCh)
	cancel() //kafka reader context
	return nil
}

func (s *StatsServer) configure() error {
	statsRepo := stats.NewMemoryStatsRepository()
	consumerManager, err := kafka.NewConsumerManager(s.config.KafkaBrokers, s.config.KafkaGroup, s.config.KafkaTopic, s.logger)
	if err != nil {
		return err
	}

	s.service = NewStatsServiceDefault(statsRepo)
	s.kafkaReader = NewKafkaReader(s.service, consumerManager, s.logger)
	s.logger.Debug("Connected to kafka: " + s.config.KafkaBrokers[0])
	return nil
}
