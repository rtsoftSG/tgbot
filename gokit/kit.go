package gokit

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"strings"
	"time"
)

type tgSDK interface {
	Send(t time.Time, lvl string, msg string) error
}

type Option func(*tgLogger)

// WithAllowLevels - define set of levels that will be send to telegram bot.
// By default send all levels.
func WithAllowLevels(levels ...level.Value) Option {
	return func(t *tgLogger) {
		t.allowLevels = make(map[level.Value]struct{}, len(levels))
		for _, lvl := range levels {
			t.allowLevels[lvl] = struct{}{}
		}
	}
}

type tgLogger struct {
	next        log.Logger
	sdk         tgSDK
	allowLevels map[level.Value]struct{}
}

//NewTgLogger - create new telegram logger.
func NewTgLogger(logger log.Logger, sdk tgSDK, opts ...Option) *tgLogger {
	lg := &tgLogger{
		next: logger,
		sdk:  sdk,
	}

	for _, opt := range opts {
		opt(lg)
	}

	return lg
}

func (l *tgLogger) Log(keyvals ...interface{}) error {
	lvl, exists := extractLVLValue(keyvals...)
	if !exists {
		return l.next.Log(keyvals...)
	}

	if l.allowLevels != nil {
		if _, exists = l.allowLevels[lvl]; !exists {
			return l.next.Log(keyvals...)
		}
	}

	if err := l.sdk.Send(time.Now(), lvl.String(), makeMessage(keyvals...)); err != nil {
		return err
	}

	return l.next.Log(keyvals...)
}

func extractLVLValue(keyvals ...interface{}) (level.Value, bool) {
	if keyvals[0] == level.Key() {
		return keyvals[1].(level.Value), true
	}
	return nil, false
}

func makeMessage(keyvals ...interface{}) string {
	msg := strings.Builder{}

	for i := 2; i < len(keyvals); i++ {
		var msgPart string
		switch v := keyvals[i].(type) {
		case string:
			msgPart = v
		case error:
			msgPart = v.Error()
		default:
			continue
		}
		msg.WriteString(msgPart)
		if i%2 == 0 {
			msg.WriteString(": ")
		} else {
			msg.WriteString("; ")
		}
	}
	return strings.TrimSuffix(msg.String(), " ")
}
