package queueconnector_test

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	queueconnector "github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/queue_connector"
	"github.com/coffeemakingtoaster/water-bottler/upload-service/pkg/util"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
)

var TEST_JOB = queueconnector.Job{
	UserEmail:   "email",
	ImageId:     "123",
	RequestTime: time.Now(),
}

func checkForMessage(connectionURI, queueName string) (queueconnector.Job, error) {
	conn, err := amqp.Dial(connectionURI)
	if err != nil {
		return queueconnector.Job{}, err
	}
	ch, err := conn.Channel()
	if err != nil {
		return queueconnector.Job{}, err
	}
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return queueconnector.Job{}, err
	}
	// read last msg
	msg := <-msgs
	var receivedJob queueconnector.Job
	err = json.Unmarshal(msg.Body, &receivedJob)
	if err != nil {
		return queueconnector.Job{}, err
	}
	return receivedJob, nil
}

func Test_publishWithoutConnection(t *testing.T) {
	queueconnector.QueueUrl = fmt.Sprintf("amqp://water:bottler@localhost:%d/", util.GetAvailablePort())

	success := queueconnector.AddJobToQueue(TEST_JOB)

	if success {
		t.Error("Expected to fail but succeeded")
	}
}

func Test_publishWithInvalidConnectionString(t *testing.T) {
	queueconnector.QueueUrl = "not valid"

	success := queueconnector.AddJobToQueue(TEST_JOB)

	if success {
		t.Error("Expected to fail but succeeded")
	}
}

func Test_publishJob(t *testing.T) {
	ctx := context.Background()
	rabbitContainer, err := rabbitmq.Run(ctx, "rabbitmq:3-management-alpine")
	if err != nil {
		t.Errorf("Could not start container due to an error %s", err.Error())
	}

	testJobInstance := TEST_JOB
	testJobInstance.RequestTime = time.Now()

	defer func() {
		if err := rabbitContainer.Terminate(context.TODO()); err != nil {
			t.Errorf("failed to terminate container: %s", err)
		}
	}()

	queueconnector.QueueUrl, err = rabbitContainer.AmqpURL(ctx)

	if err != nil {
		t.Errorf("Could not retrieve connection uri to rabbitmq due to an error: %s", err.Error())
	}

	success := queueconnector.AddJobToQueue(testJobInstance)

	if !success {
		t.Error("Expected to succeed but failed")
	}

	receivedJob, err := checkForMessage(queueconnector.QueueUrl, queueconnector.QUEUE_CHANNEL_NAME)

	if err != nil {
		t.Errorf("Could not check for received jobs due to an error: %s", err.Error())
	}

	// Funky comparison due to timestamp weirdness
	if !(receivedJob.UserEmail == testJobInstance.UserEmail &&
		receivedJob.ImageId == testJobInstance.ImageId &&
		receivedJob.RequestTime.UTC().Equal(testJobInstance.RequestTime.UTC())) {
		t.Errorf("Message received from queue did not match expected. Expected: %v\tGot: %v", TEST_JOB, receivedJob)
	}
}
