package routes

import (
	"bytes"
	"fmt"
	"io"

	"encoding/json"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/rs/zerolog/log"
	"macinvoice/internal/csv"
	"macinvoice/internal/helloservice"
	"macinvoice/internal/model"
	"macinvoice/internal/notification"
)

type Dependencies struct {
	HelloService        helloservice.HelloService
	NotificationService notification.Notification
	CSVService          csv.CSV
}

var (
	_dependencies     Dependencies
	_supportedServers map[string]string
)

// Set is used to set the routes
func Set(webserviceEngine *gin.Engine, dependencies Dependencies, supportedServers map[string]string) {
	_dependencies = dependencies
	_supportedServers = supportedServers

	// TODO: might need to consider adding default request and response timeout
	webserviceEngine.GET("/", helloWorldHandler)
	webserviceEngine.POST("/macinvoice/downloadandparse", downloadAndParseHandler)
	webserviceEngine.POST("/macinvoice/uploadandparse", uploadAndParseHandler)
}

func uploadAndParseHandler(c *gin.Context) {
	err := c.Request.ParseMultipartForm(10 << 20) // max size 10MB
	if err != nil {
		writeResponse(c, fmt.Sprintf("error while parsing multipartform: %+v", err), nil, http.StatusBadRequest)
		return
	}

	server := c.Request.FormValue("server")
	if _, isSupported := _supportedServers[strings.ToLower(server)]; !isSupported {
		writeResponse(c, fmt.Sprintf("unknown server name: %s", server), nil, http.StatusBadRequest)
		return
	}

	file, fileHeader, err := c.Request.FormFile("file")
	if err != nil {
		writeResponse(c, fmt.Sprintf("error while parsing form file: %+v", err), nil, http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		writeResponse(c, fmt.Sprintf("error while copying form file to a buffer: %+v", err), nil, http.StatusBadRequest)
		return
	}

	log.Debug().Msg(fileHeader.Filename)
	log.Debug().Msg(server)

	req := model.RequestPayload{
		Name:          server,
		URL:           "",
		Cookie:        "",
		Authorization: "",
	}
	errs := _dependencies.CSVService.Handle(req, false, buf.Bytes())
	if errs != nil && len(errs) > 0 {
		writeResponse(c, fmt.Sprintf("errors found during processing"), errs, http.StatusFailedDependency)
		return
	}

	writeResponse(c, "success", nil, http.StatusOK)

	// for debugging purposes; uncomment these lines below and comment-out writeResponse(c, "success", http.StatusOK) above
	//c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileHeader.Filename))
	//c.Data(http.StatusOK, "text/csv", buf.Bytes())

}

func downloadAndParseHandler(c *gin.Context) {

	var req model.RequestPayload
	err := json.NewDecoder(c.Request.Body).Decode(&req)
	if err != nil {
		writeResponse(c, fmt.Sprintf("error while parsing the request body: %+v", err), nil, http.StatusBadRequest)
		return
	}

	if _, isSupported := _supportedServers[strings.ToLower(req.Name)]; !isSupported {
		writeResponse(c, fmt.Sprintf("unknown server name: %s", req.Name), nil, http.StatusBadRequest)
		return
	}

	errs := _dependencies.CSVService.Handle(req, true, nil)
	if err != nil && len(errs) > 0 {
		writeResponse(c, fmt.Sprintf("errors found during processing"), errs, http.StatusFailedDependency)
		return
	}

	writeResponse(c, "success", nil, http.StatusOK)

}

func helloWorldHandler(c *gin.Context) {

	queryParams := c.Request.URL.Query()
	message := queryParams.Get("message")
	//messageValues := queryParams["message"]
	//for _, m := range messageValues {
	//	if m != "" {
	//		message = m
	//	}
	//	break
	//}

	hello := _dependencies.HelloService.WriteMessage(message)

	mObj, err := json.Marshal(hello)
	if err != nil {
		log.Err(err)
	}

	c.Header("Content-Type", binding.MIMEJSON)
	c.JSON(http.StatusOK, gin.H{
		"code": http.StatusOK,
		"obj":  string(mObj),
	})

	//c.JSON(http.StatusOK, gin.H{
	//	"message": "Hello World!",
	//})
}

func writeResponse(c *gin.Context, message string, errors []error, responseStatusCode int) {

	c.Header("Content-Type", binding.MIMEJSON)
	c.JSON(responseStatusCode, gin.H{
		"message": message,
		"errors":  errors,
	})

}
