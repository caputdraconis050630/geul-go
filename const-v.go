package main

import "time"

const (
	maxConcurrent = 5
	httpTimeout   = 30 * time.Second
	maxRetries    = 3
	retryDelay    = 5 * time.Second
)
