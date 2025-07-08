package logger

import (
	"bytes"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewZerologAdapter(t *testing.T) {
	log := NewZerologAdapter()
	assert.NotNil(t, log)
}

func TestNewWithLogger(t *testing.T) {
	log := NewZerologAdapter()
	assert.NotNil(t, log)

	newLog := NewWithLogger(log.logger)
	assert.NotNil(t, newLog)
	assert.Equal(t, log.logger, newLog.logger)
}

func TestWithField(t *testing.T) {
	buf := bytes.Buffer{}
	zlog := zerolog.New(&buf)

	log := NewWithLogger(zlog)
	assert.NotNil(t, log)

	fieldKey := "testKey"
	fieldValue := "testValue"

	log.WithField(fieldKey, fieldValue).Info("test")
	assert.JSONEq(t, buf.String(), `{"level":"info","testKey":"testValue","message":"test"}`)
}

func TestWithFields(t *testing.T) {
	buf := bytes.Buffer{}
	zlog := zerolog.New(&buf)

	log := NewWithLogger(zlog)
	assert.NotNil(t, log)

	fields := map[string]any{
		"field1": "value1",
		"field2": 42,
	}

	log.WithFields(fields).Info("test")
	assert.JSONEq(t, buf.String(), `{"level":"info","field1":"value1","field2":42,"message":"test"}`)
}

func TestInfo(t *testing.T) {
	buf := bytes.Buffer{}
	zlog := zerolog.New(&buf)

	log := NewWithLogger(zlog)
	assert.NotNil(t, log)

	log.Info("test info")
	assert.JSONEq(t, buf.String(), `{"level":"info","message":"test info"}`)
}

func TestError(t *testing.T) {
	buf := bytes.Buffer{}
	zlog := zerolog.New(&buf)

	log := NewWithLogger(zlog)
	assert.NotNil(t, log)

	log.Error("test error")
	assert.JSONEq(t, buf.String(), `{"level":"error","message":"test error"}`)
}

func TestDebug(t *testing.T) {
	buf := bytes.Buffer{}
	zlog := zerolog.New(&buf)

	log := NewWithLogger(zlog)
	assert.NotNil(t, log)

	log.Debug("test debug")
	assert.JSONEq(t, buf.String(), `{"level":"debug","message":"test debug"}`)
}

func TestWarn(t *testing.T) {
	buf := bytes.Buffer{}
	zlog := zerolog.New(&buf)

	log := NewWithLogger(zlog)
	assert.NotNil(t, log)

	log.Warn("test debug")
	assert.JSONEq(t, buf.String(), `{"level":"warn","message":"test debug"}`)
}
