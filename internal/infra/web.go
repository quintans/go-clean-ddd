package infra

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/quintans/go-clean-ddd/internal/infra/controller/web"
)

func StartWebServer(ctx context.Context, wg *sync.WaitGroup, c web.CustomerController, r web.RegistrationController) {

	e := echo.New()

	e.POST("/registrations", r.AddRegistration)
	e.GET("/customers", c.ListRegistrations)
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
		wg.Add(1)
		defer wg.Done()
		if err := e.Start(":8080"); err != nil {
			log.Fatal(err)
		} else {
			log.Println("shutting down the web server")
		}
	}()
}
