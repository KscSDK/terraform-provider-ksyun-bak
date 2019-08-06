package ksyun

import "time"

const (
	// defaultMaxRetries is default max retry attempts number
	defaultMaxRetries = 3

	// defaultInSecure is a default value to enable https
	defaultInSecure = false

	// defaultWaitInterval is the inteval to wait for state changed after resource is created
	defaultWaitInterval = 10 * time.Second

	// defaultWaitMaxAttempts is the max attempts number to wait for state changed after resource is created
	defaultWaitMaxAttempts = 10

	// defaultWaitIgnoreError is if it will ignore error during wait for state changed after resource is created
	defaultWaitIgnoreError = false

	// defaultBaseURL is the api endpoint for advanced usage
	defaultBaseURL = "https://api.ksyun.com"

	// defaultTag is the default tag for all of resources
	defaultTag = "Default"
)
const (
	// statusPending is the general status when remote resource is not completed
	statusPending = "pending"

	// statusInitialized is the general status when remote resource is completed
	statusInitialized = "initialized"

	// statusRunning is the general status when remote resource is running
	statusRunning = "running"

	// statusStopped is the general status when remote resource is stopped
	statusStopped = "stopped"
)
// trove front
const (
	tActiveStatus = "ACTIVE"
	tDeletedStatus = "DELETED"
	tCreatingStatus = "CREATING"
	tFailedStatus = "FAILED"
	tStopedStatus ="STOPPED"
)