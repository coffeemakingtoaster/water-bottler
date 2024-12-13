package queueconnector

import "time"

type Job struct {
	ImageId     string
	UserEmail   string
	RequestTime time.Time
}

func AddJobToQueue(_ Job) bool {
	// TODO: Implement me
	return true
}
