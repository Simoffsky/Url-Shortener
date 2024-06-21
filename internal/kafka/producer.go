package kafka

import (
	"context"
	"encoding/json"
	"fmt"
	"url-shorter/internal/models"
	"url-shorter/pkg/log"

	"github.com/IBM/sarama"
)

type ProducerManager struct {
	producer   sarama.AsyncProducer
	statsChan  chan models.LinkStatVisitor
	statsTopic string
	logger     log.Logger
}

func NewProducerManager(brokers []string, statTopic string, logger log.Logger) (ProducerManager, error) {
	cfg := sarama.NewConfig()
	cfg.Producer.Return.Successes = true
	producer, err := sarama.NewAsyncProducer(brokers, cfg)
	if err != nil {
		return ProducerManager{}, err
	}
	producerManager := ProducerManager{
		statsChan:  make(chan models.LinkStatVisitor),
		statsTopic: statTopic,
		producer:   producer,
		logger:     logger,
	}

	return producerManager, err
}
func (pm *ProducerManager) C() chan<- models.LinkStatVisitor {
	return pm.statsChan
}

func (pm *ProducerManager) RunProducer(ctx context.Context) {

	go func() {
		for sended := range pm.producer.Successes() {
			pm.logger.Debug("[kafka]message sended: " + fmt.Sprint(sended.Key))
		}
	}()

	go func() {
		for err := range pm.producer.Errors() {
			pm.logger.Error("[kafka]failed to send message: " + err.Error())
		}
	}()

	go func() {

		for {
			select {
			case <-ctx.Done():
				return

			case stat := <-pm.statsChan:
				pm.logger.Debug("[kafka]sending message: " + stat.LinkShort)
				jsonMsg, err := json.Marshal(stat)
				if err != nil {
					pm.logger.Debug("[kafka]failed to marshal message: " + err.Error())
					continue
				}
				pm.producer.Input() <- &sarama.ProducerMessage{
					Topic: pm.statsTopic,
					Key:   sarama.StringEncoder(stat.LinkShort),
					Value: sarama.StringEncoder(jsonMsg),
				}
			}
		}
	}()
}
