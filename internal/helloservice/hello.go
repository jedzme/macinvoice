package helloservice

import (
	"github.com/rs/zerolog/log"
	"macinvoice/internal/model"
)

type HelloService interface {
	WriteMessage(message string) model.Hello
}

type service struct {
	Config
}

func NewService(config Config) (HelloService, error) {

	return &service{config}, nil
}

// TODO: ctx is passed for traceability, we need to retrieve the transaction ID
func (s *service) WriteMessage(message string) model.Hello {
	hello := model.Hello{message}

	log.Info().Msg("WriteMessage() executed")
	log.Debug().Msg("message: " + message)

	return hello
}
