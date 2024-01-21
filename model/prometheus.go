package model

import "time"

type Alert struct {
	Status      string            `json:"status"`
	Labels      map[string]string `json:"labels"`
	Annotations map[string]string `json:annotations`
	StartsAt    time.Time         `json:"startsAt"`
	EndsAt      time.Time         `json:"endsAt"`
	StartTime   string            `json:"startTime"`
	EndTime     string            `json:"endTime"`
	Fingerprint string            `json:"fingerprint"`
	Count       int               `json:count`
}

type Notification struct {
	Version      string            `json:"version"`
	GroupKey     string            `json:"groupKey"`
	Status       string            `json:"status"`
	Receiver     string            `json:receiver`
	GroupLabels  map[string]string `json:groupLabels`
	CommonLabels map[string]string `json:commonLabels`
	ExternalURL  string            `json:externalURL`
	Alerts       []Alert           `json:alerts`
}
