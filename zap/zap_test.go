package zap

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest"
	"go.uber.org/zap/zaptest/observer"
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

func makeLoggerAndObs() (*zap.Logger, *observer.ObservedLogs, *zaptest.ShortWriter) {
	errOut := &zaptest.ShortWriter{}

	c, obs := observer.New(zap.DebugLevel)
	logger := zap.New(c, zap.ErrorOutput(errOut))
	return logger, obs, errOut
}

func TestTgLoggerLog(t *testing.T) {
	sdkMock := new(SDKMock)
	sdkMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	logger, obs, _ := makeLoggerAndObs()
	logger.WithOptions(Tg(sdkMock)).Info("msg: some message;")

	sdkMock.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"), "info", "msg: some message;")
	assert.Equal(
		t,
		zapcore.Entry{Message: "msg: some message;", Level: zapcore.InfoLevel},
		obs.AllUntimed()[0].Entry,
	)
}

func TestTgKitLoggerReturnSDKError(t *testing.T) {
	sdkMock := new(SDKMock)
	sdkMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("some error"))

	logger, _, errSync := makeLoggerAndObs()
	logger.WithOptions(Tg(sdkMock)).Error("msg: some message;")

	sdkMock.AssertCalled(t, "Send", mock.AnythingOfType("time.Time"), "error", "msg: some message;")
	assert.True(t, errSync.Called())
}

func TestTgKitLoggerFilterLogLevel(t *testing.T) {
	sdkMock := new(SDKMock)
	sdkMock.On("Send", mock.Anything, mock.Anything, mock.Anything).Return(nil)

	logger, _, _ := makeLoggerAndObs()
	logger = logger.WithOptions(TgLevels(sdkMock, zapcore.InfoLevel))

	logger.Info("some msg")
	logger.Error("some msg")

	sdkMock.AssertNumberOfCalls(t, "Send", 1)
}
