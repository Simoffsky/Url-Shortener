package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"url-shorter/internal/models"
	"url-shorter/pkg/log"

	"github.com/IBM/sarama"
)

type ConsumerManager struct {
	consumer   sarama.ConsumerGroup
	statsChan  chan models.LinkStatVisitor
	statsTopic string
	logger     log.Logger
}

func NewConsumerManager(brokers []string, groupID string, topic string, logger log.Logger) (*ConsumerManager, error) {
	config := sarama.NewConfig()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup(brokers, groupID, config)
	if err != nil {
		return nil, err
	}
	cm := &ConsumerManager{
		consumer:   consumer,
		statsChan:  make(chan models.LinkStatVisitor),
		statsTopic: topic,
		logger:     logger,
	}
	return cm, nil

}

func (cm *ConsumerManager) RunConsumer(ctx context.Context) {
	go func() {
		for err := range cm.consumer.Errors() {
			cm.logger.Error("[kafka]error while consume message: " + err.Error())
		}
	}()
	go func() {
		for {
			err := cm.consumer.Consume(ctx, []string{cm.statsTopic}, ConsumerGroupHandler{msgChan: cm.statsChan})
			if err != nil {
				cm.logger.Error("[kafka]failed to consume message: " + err.Error())
			}
		}
	}()
}

func (cm *ConsumerManager) C() <-chan models.LinkStatVisitor {
	return cm.statsChan
}

type ConsumerGroupHandler struct {
	msgChan chan models.LinkStatVisitor
}

func (ConsumerGroupHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerGroupHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h ConsumerGroupHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var message models.LinkStatVisitor
		if err := json.Unmarshal(msg.Value, &message); err != nil {
			return errors.Join(err, errors.New("failed to unmarshal message"))
		}
		h.msgChan <- message
		sess.MarkMessage(msg, "")
	}
	return nil
}
