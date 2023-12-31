package main

import (
	"flag"
	"github.com/SriHemnath/go-rest-api/src/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func init() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixNano
	debug := flag.Bool("debug", true, "sets log level to debug")

	flag.Parse()

	// Default level for this example is info, unless debug flag is present
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	//log.Debug().Msg("This message appears only when log level set to Debug")
	//log.Info().Msg("This message appears when log level set to Debug or Info")

	//if e := log.Debug(); e.Enabled() {
	//	// Compute log output only if enabled.
	//	value := "bar"
	//	e.Str("foo", value).Msg("some debug message")
	//}
}

func main() {
	log.Info().Msg("Starting up server...")
	r := server.NewServer()
	server.Start(r)
}
