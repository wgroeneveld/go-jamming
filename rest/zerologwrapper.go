package rest

import "github.com/rs/zerolog/log"

type zeroLogWrapper struct{}

func (l *zeroLogWrapper) Error(msg string, keysAndValues ...interface{}) {
	log.Error().Msgf(msg, keysAndValues...)
}
func (l *zeroLogWrapper) Info(msg string, keysAndValues ...interface{}) {
	log.Info().Msgf(msg, keysAndValues...)
}
func (l *zeroLogWrapper) Debug(msg string, keysAndValues ...interface{}) {
	// no thanks.
	// log.Debug().Msgf(msg, keysAndValues...)
}
func (l *zeroLogWrapper) Warn(msg string, keysAndValues ...interface{}) {
	log.Warn().Msgf(msg, keysAndValues...)
}
