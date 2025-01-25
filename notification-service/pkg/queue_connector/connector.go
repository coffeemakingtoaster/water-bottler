package queueconnector

import (
	"context"
	"encoding/json"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rs/zerolog/log"
)

var queueConn *amqp.Connection
var queue *amqp.Queue
var queueChan *amqp.Channel
var QueueUrl string
var QUEUE_URL_ENV = "QUEUE_URL"
var QUEUE_CHANNEL_NAME = "finished-jobs"

type FinishedJob struct {
	ImageId   string `json:"image_id"`
	UserEmail string `json:"email"`
}

func init() {
	QueueUrl = os.Getenv(QUEUE_URL_ENV)
	if len(QueueUrl) == 0 {
		QueueUrl = "amqp://water:bottler@localhost:5672/"
	}
}

func getChannel() (*amqp.Channel, *amqp.Queue) {
	if queueChan != nil && queue != nil {
		if queueChan.IsClosed() || queueConn.IsClosed() {
			log.Info().Msg("Queue connection or channel have been closed. Restablishing connection...")
			queueChan = nil
			queueConn = nil
			return getChannel()
		}
		return queueChan, queue
	}
	if queueConn == nil {
		new_conn, err := amqp.Dial(QueueUrl)
		if err != nil {
			log.Warn().Msgf("Could not contact queue due to an error: %s", err.Error())
			return nil, nil
		}
		queueConn = new_conn
	}

	var err error
	queueChan, err = queueConn.Channel()
	if err != nil {
		log.Warn().Msgf("Could not create channel dur to an error: %s", err.Error())
		queueChan = nil
		return nil, nil
	}

	new_queue, err := queueChan.QueueDeclare(
		QUEUE_CHANNEL_NAME,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Warn().Msgf("Could not create queue dur to an error: %s", err.Error())
		queueChan = nil
		return nil, nil
	}
	queue = &new_queue

	return queueChan, queue
}

func AddJobToQueue(job FinishedJob) bool {
	jobParsed, err := json.Marshal(job)

	if err != nil {
		log.Warn().Msgf("Could not marshal job body due to an error: %s", err.Error())
		return false
	}

	ch, q := getChannel()
	if ch == nil || q == nil {
		return false
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(ctx,
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         jobParsed,
		})

	if err != nil {
		log.Warn().Msgf("Could not publish message to queue to an error: %s", err.Error())
		return false
	}

	return true
}

func ConsumeJobFromQueue() (chan FinishedJob, bool) {
	ch, q := getChannel()
	if ch == nil || q == nil {
		log.Error().Msg("Could not get channel or queue")
		return nil, false
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		log.Error().Err(err).Msg("Could not consume message from queue")
		return nil, false
	}

	finishedJobs := make(chan FinishedJob)
	go func() {
		for msg := range msgs {
			var job FinishedJob
			err := json.Unmarshal(msg.Body, &job)
			if err != nil {
				log.Warn().Msgf("Could not unmarshal job due to an error: %s", err.Error())
				continue
			}
			finishedJobs <- job
		}
	}()

	return finishedJobs, true
}
