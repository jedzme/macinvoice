package main

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"macinvoice/cmd/webapp/config"
	"macinvoice/cmd/webapp/routes"
	"macinvoice/internal/csv"
	"macinvoice/internal/helloservice"
	mHTTP "macinvoice/internal/http"
	"macinvoice/internal/notification"
)

func main() {

	// init config
	port := config.AppConfig.GetString("SERVER_PORT")

	supportedServers := config.AppConfig.GetStringMapString("SUPPORTED_SERVERS")
	if len(supportedServers) == 0 {
		log.Fatal().Msg("Please check SUPPORTED_SERVERS.")
	}

	// init dependencies
	dependencies, err := getDependencies()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to initialize dependencies: ")
	}

	// init the webservice
	webserviceEngine := gin.New()
	if config.AppConfig.GetBool("PROD_MODE") {
		gin.SetMode(gin.ReleaseMode)
	}
	routes.Set(webserviceEngine, dependencies, supportedServers)

	// TODO: might need to create a logger interface for the project that could accept the Context from the webservice request for traceability
	// TODO: might need to SetTrustedProxies()
	log.Info().Msg("Starting application on port: " + port)
	webserviceEngine.Run(":" + port)

}

// For configurations, use the config.AppConfig right away inside this method
func getDependencies() (routes.Dependencies, error) {

	dependencies := routes.Dependencies{}

	//  === Service A ===
	helloConfig := helloservice.Config{
		A: "",
		B: false,
	}
	helloService, err := helloservice.NewService(helloConfig)
	if err != nil {
		return routes.Dependencies{}, err
	}
	dependencies.HelloService = helloService

	// == Notification Service ===
	notificationConfig := notification.Config{
		A: 0,
		B: false,
	}
	notificationSvc, err := notification.NewService(notificationConfig)
	if err != nil {
		return routes.Dependencies{}, err
	}
	dependencies.NotificationService = notificationSvc

	// == HTTP Service ===
	httpServiceConfig := mHTTP.Config{
		MaxRetries:    config.AppConfig.GetInt("HTTP_CALLS_MAX_RETRIES"),
		ClientTimeout: config.AppConfig.GetDuration("HTTP_CLIENT_TIMEOUT"),
	}
	restService, err := mHTTP.NewService(httpServiceConfig)
	if err != nil {
		return routes.Dependencies{}, err
	}

	// == CSV Parser Service ===
	csvConfig := csv.Config{
		A: 0,
		B: false,
	}
	csvService, err := csv.NewService(csvConfig, restService)
	if err != nil {
		return routes.Dependencies{}, err
	}
	dependencies.CSVService = csvService

	return dependencies, nil

}
