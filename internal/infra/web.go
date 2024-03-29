package infra

import (
	"context"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
	"github.com/quintans/toolkit/latch"
)

func StartWebServer(
	ctx context.Context,
	lock *latch.CountDownLatch,
	cfg WebConfig,
	c web.CustomerController,
	r web.RegistrationController,
) {
	e := echo.New()

	e.POST("/registrations", r.AddRegistration)
	e.GET("/registrations/:id", r.ConfirmRegistration)
	e.GET("/customers", c.ListCustomers)
	e.PATCH("/customers/:id", c.UpdateCustomer)

	go func() {
		<-ctx.Done()
		c, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(c); err != nil {
			log.Printf("[ERROR] %+v", err)
		}
	}()

	lock.Add(1)
	go func() {
		defer lock.Done()
		err := e.Start(cfg.Port)
		if err == nil || errors.Is(err, http.ErrServerClosed) {
			log.Println("web server closed")
		} else {
			log.Printf("[ERROR] %+v", err)
		}
	}()
}
