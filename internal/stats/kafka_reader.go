package stats

import (
	"context"

	"url-shorter/internal/kafka"
	"url-shorter/pkg/log"
)

type KafkaReader struct {
	service      StatsService
	kafkaManager *kafka.ConsumerManager
	logger       log.Logger
}

func NewKafkaReader(service StatsService, kafkaManager *kafka.ConsumerManager, logger log.Logger) *KafkaReader {
	return &KafkaReader{
		service:      service,
		kafkaManager: kafkaManager,
		logger:       logger,
	}
}

func (k *KafkaReader) Start(ctx context.Context) {
	k.kafkaManager.RunConsumer(ctx)
	go func() {
		for stat := range k.kafkaManager.C() {
			err := k.service.SendStat(&stat)
			if err != nil {
				k.logger.Error("Error while creating entity in repository:" + err.Error())
			}
			k.logger.Debug("added stat: " + stat.LinkShort)
		}
	}()
}
