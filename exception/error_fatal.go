package exception

import (
	"github.com/rs/zerolog/log"
)

func FatalIfNeeded(err error, msg string) {
	if err != nil {
		log.Fatal().Err(err).Msg(msg)
	}
}
