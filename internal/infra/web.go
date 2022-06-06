package infra

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
	"github.com/quintans/toolkit/latch"
)

func StartWebServer(
	ctx context.Context,
	lock *latch.CountDownLatch,
	address string,
	c web.CustomerController,
	r web.RegistrationController,
) {
	e := echo.New()

	e.GET("/registrations/:id", r.ConfirmRegistration)
	e.POST("/registrations", r.AddRegistration)
	e.GET("/customers", c.ListCustomers)
	e.PATCH("/customers/:id", c.UpdateRegistration)

	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(c); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		lock.Add(1)
		defer lock.Done()
		if err := e.Start(address); err != nil {
			log.Fatal(err)
		} else {
			log.Println("shutting down the web server")
		}
	}()
}
