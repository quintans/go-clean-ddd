package main

import (
	"log"

	"github.com/quintans/go-clean-ddd/internal/controller/web"
	"github.com/quintans/go-clean-ddd/internal/gateway/postgres"
	pg "github.com/quintans/go-clean-ddd/internal/infra/postgres"
	iWeb "github.com/quintans/go-clean-ddd/internal/infra/web"
	"github.com/quintans/go-clean-ddd/internal/usecase/command"
	"github.com/quintans/go-clean-ddd/internal/usecase/query"
)

func main() {
	db, err := pg.New()
	if err != nil {
		log.Fatal(err)
	}

	repo := postgres.NewCustomerRepository(db)

	commands := command.NewCommander(repo, repo)
	queries := query.NewQuerier(repo)

	c := web.NewController(commands, queries)

	if err := iWeb.StartWebServer(c); err != nil {
		log.Fatal(err)
	}
}
