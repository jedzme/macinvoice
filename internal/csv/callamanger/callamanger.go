package callamanger

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/rs/zerolog/log"
	"io"
	mhttp "macinvoice/internal/http"
	"macinvoice/internal/model"
	"net/http"
	"strings"
)

func Process(restService mhttp.REST, downloadFirst bool, url string, cookie string, authorization string, csvBytes []byte) []error {
	log.Info().Msg("callamanger.Process() executed")

	errs := make([]error, 0)
	records := make(map[string][]model.Record)
	var err error
	if downloadFirst {
		csvBytes, err = download(restService, url, cookie, authorization)
		if err != nil {
			return append(errs, err)
		}
	}

	csvReader := csv.NewReader(bytes.NewBuffer(csvBytes))
	var totalRecordsCounter int32

	isFirstLine := true
	for {

		line, readErr := csvReader.Read()
		if readErr != nil {
			if readErr == io.EOF || errors.Is(readErr, io.EOF) {
				log.Info().Msg("EOF reached.")
			} else {
				log.Error().Err(readErr)
				errs = append(errs, readErr)
			}
			break
		}

		if isFirstLine {
			isFirstLine = false
			continue
		}

		totalRecordsCounter++

		recordStruct, cErr := createRecordStruct(line)
		if cErr != nil {
			log.Error().Err(cErr)
			errs = append(errs, cErr)
			continue
		} else {
			sErr := seggregateRecordStruct(recordStruct, records)
			if sErr != nil {
				errs = append(errs, sErr)
				continue
			}
		}
	}

	log.Info().Msg(fmt.Sprintf("Total Records Counter: %d", totalRecordsCounter))

	for company, companyRecords := range records {
		log.Info().Msg(fmt.Sprintf("Company: %s; Total Records: %d", company, len(companyRecords)))

	}

	return errs
}

func seggregateRecordStruct(record model.Record, recordMap map[string][]model.Record) error {
	company := record.Company
	if company == "" {
		return fmt.Errorf("company name is empty: %+v", record)
	}

	companyRecords := recordMap[company]
	if companyRecords == nil {
		companyRecords = make([]model.Record, 0)
	}
	companyRecords = append(companyRecords, record)

	// update the companyRecords
	recordMap[company] = companyRecords

	return nil
}

func createRecordStruct(rawCSVLine []string) (model.Record, error) {

	if len(rawCSVLine) != 16 {
		lineErr := fmt.Errorf("length of csv line is not 16: %d, csv line value: %+v", len(rawCSVLine), rawCSVLine)
		return model.Record{}, lineErr
	}

	if rawCSVLine[0] == "" {
		return model.Record{}, fmt.Errorf("company name is empty: %+v", rawCSVLine)
	}

	return model.Record{
		Company:    rawCSVLine[0],
		Person:     rawCSVLine[1],
		Name:       rawCSVLine[2],
		DeviceType: rawCSVLine[3],
		MacAddress: rawCSVLine[4],
		Registered: func() bool {
			if strings.EqualFold(rawCSVLine[5], "true") {
				return true
			}
			return false

		}(),
		Status:           rawCSVLine[6],
		UUIDCreationDate: rawCSVLine[7],
		DownloadDate:     rawCSVLine[8],
		HotDesking: func() bool {
			if strings.EqualFold(rawCSVLine[9], "true") {
				return true
			}
			return false

		}(),
		HotDeskingID:    rawCSVLine[10],
		HotDeskingPhone: rawCSVLine[11],
		Location:        rawCSVLine[12],
		Group:           rawCSVLine[13],
		Comment:         rawCSVLine[14],
		Firmware:        rawCSVLine[15],
	}, nil

}

func download(restService mhttp.REST, url string, cookie string, authorization string) ([]byte, error) {
	log.Info().Msg("callamanger.download() executed")
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
