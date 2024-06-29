package observability

import (
	"context"
	"errors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

func Init(ctx context.Context, cfg Cfg, serviceName, namespace string) func() {
	// log
	if cfg.LogPretty {
		zerolog.LevelColors[zerolog.DebugLevel] = 35
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Caller().Logger()

	} else {
		log.Logger = log.With().Str("service", serviceName).Caller().Logger()
	}
	ctx = log.Logger.WithContext(ctx)

	level, err := zerolog.ParseLevel(cfg.LogLevel)
	if err != nil {
		panic(err)
	}
	zerolog.SetGlobalLevel(level)
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	// trace
	if cfg.Trace {

		otelShutdown, err := SetupOTelSDK(ctx, cfg.OltpPass.GetSecret(), cfg.OltpEndpoint, serviceName, namespace)
		if err != nil {
			log.Fatal().Err(err).Msg("observability SetupOTelSDK")
			return nil
		}
		// Handle shutdown properly so nothing leaks.
		f := func() {
			err = errors.Join(err, otelShutdown(context.Background()))
			if err != nil {
				log.Error().Err(err).Msg("otelShutdown")
			}
		}
		return f
	}
	return nil
}
