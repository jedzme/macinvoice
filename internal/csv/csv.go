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
	_callamangerProcess = func(restService mhttp.REST, downloadFirst bool, url string, cookie string, authorization string, csvBytes []byte) []error {
		return callamanger.Process(restService, downloadFirst, url, cookie, authorization, csvBytes)
	}
)

type CSV interface {
	Handle(req model.RequestPayload, downloadFirst bool, csvBytes []byte) []error
}

type service struct {
	Config
}

func NewService(config Config, restService mhttp.REST) (CSV, error) {
	_restService = restService
	return &service{config}, nil
}

func (s *service) Handle(req model.RequestPayload, downloadFirst bool, csvBytes []byte) []error {
	errs := make([]error, 0)
	server, err := model.GetServerEnum(strings.ToLower(req.Name))
	if err != nil {
		return append(errs, err)
	}

	switch server {
	case model.CALLAMANGER:
		url := req.URL
		if downloadFirst {
			if url == "" {
				return append(errs, fmt.Errorf("url is required"))
			}

			if req.Cookie == "" {
				return append(errs, fmt.Errorf("cookie is required"))
			}

			if req.Authorization == "" {
				return append(errs, fmt.Errorf("authorization is required"))
			}
		} else {
			if len(csvBytes) == 0 {
				return append(errs, fmt.Errorf("no csv file uploaded"))
			}
		}

		errs = _callamangerProcess(_restService, downloadFirst, url, req.Cookie, req.Authorization, csvBytes)
	default:
		errs = append(errs, fmt.Errorf("unknown server"))
	}

	return errs
}

// TODO: for now save to local
func saveCSV(filename string, csvBytes []byte) error {
	return nil
}
