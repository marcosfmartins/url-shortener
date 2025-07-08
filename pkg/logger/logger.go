package logger

import (
	"github.com/marcosfmartins/url_shortener/internal/entity"
	"github.com/rs/zerolog"
	"os"
)

type ZerologAdapter struct {
	logger zerolog.Logger
}

func NewZerologAdapter() *ZerologAdapter {
	base := zerolog.New(os.Stdout).With().Timestamp().Logger()
	return &ZerologAdapter{
		logger: base,
	}
}

func NewWithLogger(z zerolog.Logger) *ZerologAdapter {
	return &ZerologAdapter{
		logger: z,
	}
}

func (z *ZerologAdapter) WithField(key string, value any) entity.Logger {
	return NewWithLogger(z.logger.With().Interface(key, value).Logger())
}

func (z *ZerologAdapter) Err(err error) entity.Logger {
	return NewWithLogger(z.logger.With().Err(err).Logger())
}

func (z *ZerologAdapter) WithFields(fields map[string]any) entity.Logger {
	ctx := z.logger.With()
	for k, v := range fields {
		ctx = ctx.Interface(k, v)
	}
	return NewWithLogger(ctx.Logger())
}

func (z *ZerologAdapter) Info(msg string) {
	z.logger.Info().Msg(msg)
}

func (z *ZerologAdapter) Error(msg string) {
	z.logger.Error().Msg(msg)
}

func (z *ZerologAdapter) Debug(msg string) {
	z.logger.Debug().Msg(msg)
}

func (z *ZerologAdapter) Warn(msg string) {
	z.logger.Warn().Msg(msg)
}

func (z *ZerologAdapter) Fatal(msg string) {
	z.logger.Fatal().Msg(msg)
}
