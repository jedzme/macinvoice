package csv

import (
	"fmt"
	"macinvoice/internal/csv/callamanger"
	mhttp "macinvoice/internal/http"
	"macinvoice/internal/model"
	"strings"
)

var (
	_restService        mhttp.REST
	_callamangerProcess = func(restService mhttp.REST, downloadFirst bool, url string, cookie string, authorization string, csvBytes []byte) error {
		return callamanger.Process(restService, downloadFirst, url, cookie, authorization, csvBytes)
	}
)

type CSV interface {
	Handle(req model.RequestPayload, downloadFirst bool, csvBytes []byte) error
}

type service struct {
	Config
}

func NewService(config Config, restService mhttp.REST) (CSV, error) {
	_restService = restService
	return &service{config}, nil
}

func (s *service) Handle(req model.RequestPayload, downloadFirst bool, csvBytes []byte) error {

	server, err := model.GetServerEnum(strings.ToLower(req.Name))
	if err != nil {
		return err
	}

	switch server {
	case model.CALLAMANGER:
		url := req.URL
		if downloadFirst {
			if url == "" {
				return fmt.Errorf("url is required")
			}

			if req.Cookie == "" {
				return fmt.Errorf("cookie is required")
			}

			if req.Authorization == "" {
				return fmt.Errorf("authorization is required")
			}
		} else {
			if len(csvBytes) == 0 {
				return fmt.Errorf("no csv file uploaded")
			}
		}

		err = _callamangerProcess(_restService, downloadFirst, url, req.Cookie, req.Authorization, csvBytes)
	default:
		err = fmt.Errorf("unknown server")
	}

	return err
}
