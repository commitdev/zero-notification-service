// Helpers for ECS log formatting
// https://www.elastic.co/guide/en/ecs/master/index.html
package log

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// ECSHTTP represents a subset of the fields of the ECS HTTP object
type ECSHTTP struct {
	Request  ECSRequest
	Response ECSResponse
}

func (o ECSHTTP) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddObject("request", o.Request)
	return enc.AddObject("response", o.Response)
}

// ECSResponse represents a subset of the fields of the ECS HTTP Response object
type ECSResponse struct {
	StatusCode int
}

func (o ECSResponse) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddInt("status_code", o.StatusCode)
	return nil
}

// ECSRequest represents a subset of the fields of the ECS HTTP Request object
type ECSRequest struct {
	Method string
}

func (o ECSRequest) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("method", o.Method)
	return nil
}

// ECSURL represents a subset of the fields of the ECS URL object
type ECSURL struct {
	Original string
}

func (o ECSURL) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("original", o.Original)
	return nil
}

// ECSEvent represents a subset of the fields of the ECS Event object
type ECSEvent struct {
	Action   string
	Duration time.Duration
}

func (o ECSEvent) MarshalLogObject(enc zapcore.ObjectEncoder) error {
	enc.AddString("action", o.Action)
	enc.AddInt64("duration", int64(o.Duration))
	return nil
}
