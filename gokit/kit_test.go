package gokit

import (
	"context"
	"errors"
	"github.com/go-kit/kit/log/level"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

type SDKMock struct {
	mock.Mock
}

func (s *SDKMock) Send(_ context.Context, t time.Time, lvl string, msg string) error {
	args := s.Called(t, lvl, msg)
	return args.Error(0)
}

type kitLoggerMock struct {
	mock.Mock
}

func (k *kitLoggerMock) Log(keyvals ...interface{}) error {
	args := k.Called(keyvals)
	return args.Error(0)
}

func TestTgKitLoggerLog(t *testing.T) {
	kitLogger := new(kitLoggerMock)
	kitLogger.On("Log", mock.Anything).Return(nil)
	sdkMock := new(SDKMock)
	sdkMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	logger := NewTgLogger(kitLogger, sdkMock)

	assert.NoError(t, level.Error(logger).Log("msg", "some message"))

	sdkMock.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"), "error", "msg: some message;")
	kitLogger.AssertCalled(t, "Log", mock.Anything)
}

func TestTgKitLoggerDoesNotLogWithoutLvl(t *testing.T) {
	kitLogger := new(kitLoggerMock)
	kitLogger.On("Log", mock.Anything).Return(nil)
	sdkMock := new(SDKMock)
	sdkMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	logger := NewTgLogger(kitLogger, sdkMock)

	assert.NoError(t, logger.Log("msg", "some message"))

	sdkMock.AssertNotCalled(t, "Send", mock.Anything, mock.Anything, mock.Anything)
	kitLogger.AssertCalled(t, "Log", mock.Anything)
}

func TestTgKitLoggerReturnSDKError(t *testing.T) {
	kitLogger := new(kitLoggerMock)
	kitLogger.On("Log", mock.Anything).Return(nil)
	sdkMock := new(SDKMock)
	sdkMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error"))

	logger := NewTgLogger(kitLogger, sdkMock)

	assert.Error(t, level.Error(logger).Log("msg", "some message"))

	sdkMock.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"), "error", "msg: some message;")
	kitLogger.AssertNotCalled(t, "Log", mock.Anything)
}

func TestTgKitLoggerFilterLogLevel(t *testing.T) {
	kitLogger := new(kitLoggerMock)
	kitLogger.On("Log", mock.Anything).Return(nil)
	sdkMock := new(SDKMock)
	sdkMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	logger := NewTgLogger(kitLogger, sdkMock, WithAllowLevels(level.ErrorValue()))

	assert.NoError(t, level.Error(logger).Log("msg", "some message"))
	assert.NoError(t, level.Info(logger).Log("msg", "some message"))

	sdkMock.AssertNumberOfCalls(t, "Send", 1)
	kitLogger.AssertNumberOfCalls(t, "Log", 2)
}
