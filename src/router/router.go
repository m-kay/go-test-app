package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/m-kay/test-app/src/service"
)

type Router struct {
	app          *fiber.App
	helloService *service.HelloService
}

func NewRouter(app *fiber.App, helloService *service.HelloService) *Router {
	r := &Router{app: app, helloService: helloService}
	r.setupRoutes()
	return r
}

func (router *Router) setupRoutes() {
	router.app.Get("hello", func(c *fiber.Ctx) error {
		helloResponse, err := router.helloService.GetHello()
		if err != nil {
			c.Status(500).SendString(err.Error())
			return err
		}
		return c.JSON(Response{Message: "Hello " + helloResponse})
	})
}

type Response struct {
	Message string `json:"message"`
}
