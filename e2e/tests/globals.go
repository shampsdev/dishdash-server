package tests

import (
	"net/http"
	"time"
)

var (
	ApiHost    string
	SIOHost    string
	httpClient = &http.Client{Timeout: 10 * time.Second}
	waitTime   = 10 * time.Second
)
