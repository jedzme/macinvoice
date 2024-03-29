package model

import "fmt"

type RequestPayload struct {
	Name          string `json:"name"`
	URL           string `json:"url"`
	Cookie        string `json:"cookie"`
	Authorization string `json:"authorization"`
}

type Hello struct {
	Message string `json:"message"`
}

type SERVER string

const (
	CALLAMANGER SERVER = "callamanger"
)

type RESTResponse struct {
	Code  int
	Body  []byte
	Error error
}

func GetServerEnum(server string) (SERVER, error) {

	switch server {
	case "callamanger":
		return CALLAMANGER, nil
	}

	return "", fmt.Errorf("unknown server")
}

type Record struct {
	Company          string
	Person           string
	Name             string
	DeviceType       string
	MacAddress       string
	Registered       bool // string in csv (Yes/No)
	Status           string
	UUIDCreationDate string
	DownloadDate     string
	HotDesking       bool // string in csv (Yes/No)
	HotDeskingID     string
	HotDeskingPhone  string
	Location         string
	Group            string
	Comment          string
	Firmware         string
}
