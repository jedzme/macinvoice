package http

import (
	"fmt"

	"io/ioutil"
	"net/http"

	"github.com/rs/zerolog/log"
	"macinvoice/internal/model"
)

type REST interface {
	POST(url string, headers map[string]string, requestPayload []byte) model.RESTResponse
	GET(url string, headers map[string]string) model.RESTResponse
}

type service struct {
	Config
}

func NewService(config Config) (REST, error) {

	return &service{config}, nil
}

func (s *service) GET(url string, headers map[string]string) model.RESTResponse {

	var retries = 0
	var err error
	var respBytes []byte
	var statusCode = 0
	var resp *http.Response
	var req *http.Request

	client := http.Client{
		Transport:     nil,
		CheckRedirect: nil,
		Jar:           nil,
		Timeout:       s.ClientTimeout,
	}

	req, err = http.NewRequest(http.MethodGet, url, nil) // TODO: we might need to send a request body soon
	if err != nil {
		return model.RESTResponse{
			Code:  statusCode,
			Body:  respBytes,
			Error: err,
		}
	}

	log.Debug().Msg("==HEADERS[START]==")
	for k, v := range headers {
		log.Debug().Msg(fmt.Sprintf("k: %s; v: %s", k, v))
		req.Header.Set(k, v)
	}
	log.Debug().Msg("==HEADERS[END]==")

	for {

		if retries == s.MaxRetries {
			break
		}

		resp, err = client.Do(req)
		if err == nil && resp != nil {

			statusCode = resp.StatusCode
			respBytes, err = ioutil.ReadAll(resp.Body)

			if statusCode == http.StatusOK {
				break
			}

		}

		retries++

	}

	defer resp.Body.Close()

	return model.RESTResponse{
		Code:  statusCode,
		Body:  respBytes,
		Error: err,
	}
}

func (s *service) POST(url string, headers map[string]string, requestPayload []byte) model.RESTResponse {
	panic("implement me")
	return model.RESTResponse{}
}
