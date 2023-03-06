package callamanger

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	mhttp "macinvoice/internal/http"
	"net/http"
)

func Process(restService mhttp.REST, downloadFirst bool, url string, cookie string, authorization string, csvBytes []byte) error {
	var err error
	if downloadFirst {
		csvBytes, err = download(restService, url, cookie, authorization)
		if err != nil {
			return err
		}
	}

	csvReader := csv.NewReader(bytes.NewBuffer(csvBytes))
	for {
		line, err := csvReader.Read()
		if err != nil {
			if err == io.EOF || errors.Is(err, io.EOF) {
				log.Info().Msg("EOF reached.")
			} else {
				log.Error().Err(err)
			}
			break
		}
		log.Info().Msg(fmt.Sprintf("%+v", line))
	}

	return nil
}

func download(restService mhttp.REST, url string, cookie string, authorization string) ([]byte, error) {

	headers := make(map[string]string)
	headers["Cookie"] = cookie
	headers["Content-Type"] = "text/csv"
	headers["Authorization"] = authorization

	resp := restService.GET(url, headers)

	if resp.Error != nil {
		return nil, resp.Error
	}

	if resp.Code != http.StatusOK {
		return nil, fmt.Errorf("download returned a failed status code: %d; body: %s", resp.Code, string(resp.Body))
	}

	return resp.Body, nil
}
