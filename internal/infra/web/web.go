package web

import (
	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/controller/web"
)

func StartWebServer(c web.CustomerController, r web.RegistrationController) error {
	e := echo.New()

	e.POST("/registrations", r.AddRegistration)
	e.GET("/customers", c.ListRegistrations)
	e.PATCH("/customers/:id", c.UpdateRegistration)

	return e.Start(":8080")
}
