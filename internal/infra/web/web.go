package web

import (
	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/controller/web"
)

func StartWebServer(c web.AdminController) error {
	e := echo.New()

	e.GET("/admin/registrations", c.ListRegistrations)
	e.POST("/admin/registrations", c.AddRegistration)
	e.PATCH("/admin/registrations", c.UpdateRegistration)

	return e.Start(":8080")
}
