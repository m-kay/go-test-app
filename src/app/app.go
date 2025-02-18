package app

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/m-kay/test-app/src/router"
	"github.com/m-kay/test-app/src/service"
	"os"
	"strconv"
)

type Application struct {
	Fiber  *fiber.App
	router *router.Router
}

func NewApplication() *Application {
	fiberApp := fiber.New()
	port, ok := os.LookupEnv("HELLO_SERVICE_PORT")
	if !ok {
		port = "5000"
	}
	portNumber, _ := strconv.Atoi(port)
	helloService, err := service.NewHelloService(portNumber)
	if err != nil {
		panic(err)
	}
	r := router.NewRouter(fiberApp, helloService)
	app := Application{fiberApp, r}
	return &app
}

func (app *Application) Start(port uint) error {
	formatString := fmt.Sprintf(":%d", port)
	return app.Fiber.Listen(formatString)
}

func (app *Application) Stop() error {
	return app.Fiber.Shutdown()
}
