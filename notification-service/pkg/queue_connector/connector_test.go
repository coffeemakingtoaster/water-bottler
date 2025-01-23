package queueconnector_test

import (
	"fmt"
	"testing"

	queueconnector "github.com/coffeemakingtoaster/water-bottler/notification-service/pkg/queue_connector"
	"github.com/coffeemakingtoaster/water-bottler/notification-service/pkg/util"
)

var TEST_JOB = queueconnector.FinishedJob{
	UserEmail: "email",
	ImageId:   "123",
}

func Test_consumeWithoutConnection(t *testing.T) {
	queueconnector.QueueUrl = fmt.Sprintf("amqp://water:bottler@localhost:%d/", util.GetAvailablePort())

	_, success := queueconnector.ConsumeJobFromQueue()

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
