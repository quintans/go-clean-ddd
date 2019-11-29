package framework

import (
	"github.com/quintans/go-clean-ddd/internal/adapter/persistence"
	"github.com/quintans/go-clean-ddd/internal/adapter/web"
	"github.com/quintans/go-clean-ddd/internal/domain/service"
	"github.com/quintans/go-clean-ddd/internal/usecase"
	"github.com/labstack/echo/v4"
	"github.com/quintans/goSQL/translators"

	// initializes the postgres driver
	_ "github.com/lib/pq"
)

// NewService creates the service server
func NewService(address string) error {
	pool, err := persistence.Pool("postgres", "dbname=postgres user=postgres password=secret port=5432 sslmode=disable")
	if err != nil {
		return err
	}

	txMan := persistence.NewTxManager(pool)

	repo := persistence.NewCustomerRepositoryImpl(translators.NewPostgreSQLTranslator())
	customerSvc := service.CustomerServiceImpl{
		CustomerRepository: repo,
	}
	var svc usecase.RegistrationService = usecase.NewRegistrationServiceImpl(customerSvc, repo)
	svc = persistence.NewRegistrationServiceTx(txMan, svc)

	ctrl := web.NewAdminController(svc)

	e := route(ctrl)
	err = e.Start(address)
	e.Logger.Fatal(err)
	return err
}

func route(ctrl web.AdminController) *echo.Echo {
	e := echo.New()
	e.GET("/admin/registrations", ctrl.ListRegistrations)
	e.POST("/admin/registrations", ctrl.AddRegistration)
	return e
}
