package zap

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

type tgSDK interface {
	Send(ctx context.Context, t time.Time, lvl string, msg string) error
}

func Tg(sdk tgSDK) zap.Option {
	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			return sdk.Send(context.Background(), entry.Time.UTC(), entry.Level.String(), entry.Message)
		})
	})
}

func TgLevels(sdk tgSDK, levels ...zapcore.Level) zap.Option {
	allowLvls := make(map[zapcore.Level]struct{}, len(levels))
	for _, lvl := range levels {
		allowLvls[lvl] = struct{}{}
	}

	return zap.WrapCore(func(core zapcore.Core) zapcore.Core {
		return zapcore.RegisterHooks(core, func(entry zapcore.Entry) error {
			if _, exists := allowLvls[entry.Level]; exists {
				return sdk.Send(context.Background(), entry.Time.UTC(), entry.Level.String(), entry.Message)
			}
			return nil
		})
	})
}
