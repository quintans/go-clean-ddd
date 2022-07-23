package infra_test

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/docker/go-connections/nat"
	"github.com/jmoiron/sqlx"
	"github.com/quintans/faults"
	"github.com/quintans/go-clean-ddd/internal/infra"
	"github.com/quintans/go-clean-ddd/internal/infra/gateway/postgres"
	"github.com/quintans/toolkit/latch"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const baseUrl = "http://localhost:8080"

var db *sqlx.DB

func TestRegister(t *testing.T) {
	resp, err := http.Post(baseUrl+"/registrations", "application/json; charset=UTF-8", strings.NewReader(`{"email":"abc@xpto.pt"}`))
	require.NoError(t, err)
	defer resp.Body.Close()

	require.Equal(t, http.StatusCreated, resp.StatusCode)
	body, err := ioutil.ReadAll(resp.Body)
	require.NoError(t, err)
	require.NotEmpty(t, body)

	registrations := getAllRegistrations(t)
	require.Len(t, registrations, 1)
	reg := registrations[0]
	require.Equal(t, "abc@xpto.pt", reg.Email)
	require.Equal(t, false, reg.Verified)

	// give time for the outbox flush to trigger
	time.Sleep(2 * time.Second)

	registrations = getAllRegistrations(t)
	require.Len(t, registrations, 1)
	reg = registrations[0]
	require.Equal(t, "abc@xpto.pt", reg.Email)
	require.Equal(t, true, reg.Verified)
}

func getAllRegistrations(t *testing.T) []postgres.Registration {
	var registrations []postgres.Registration
	err := db.Select(&registrations, "SELECT * FROM registrations")
	require.NoError(t, err)
	return registrations
}

func TestMain(m *testing.M) {
	lock := latch.NewCountDownLatch()
	dbCfg, _, err := setup()
	if err != nil {
		panic(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	dbConfig := infra.DbConfig{
		DbName:     dbCfg.Database,
		DbHost:     dbCfg.Host,
		DbPort:     dbCfg.Port,
		DbUser:     dbCfg.Username,
		DbPassword: dbCfg.Password,
	}
	infra.Start(ctx, lock, infra.Config{
		OutboxHeartbeat: time.Second,
		DbConfig:        dbConfig,
		WebConfig: infra.WebConfig{
			Port: ":8080",
		},
	})

	db = infra.NewDB(dbConfig)

	code := m.Run()

	cancel()
	lock.WaitWithTimeout(3 * time.Second)

	os.Exit(code)
}

type DBConfig struct {
	Database string
	Host     string
	Port     int
	Username string
	Password string
}

func setup() (DBConfig, func(), error) {
	dbConfig := DBConfig{
		Database: "registration",
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		Password: "secret",
	}
	tcpPort := strconv.Itoa(dbConfig.Port)
	natPort := nat.Port(tcpPort)

	req := testcontainers.ContainerRequest{
		Image:        "postgres:12.3",
		ExposedPorts: []string{tcpPort + "/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     dbConfig.Username,
			"POSTGRES_PASSWORD": dbConfig.Password,
			"POSTGRES_DB":       dbConfig.Database,
		},
		WaitingFor: wait.ForListeningPort(natPort),
	}
	ctx := context.Background()
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return DBConfig{}, nil, faults.Wrap(err)
	}

	tearDown := func() {
		// this is crashing without a trace
		container.Terminate(context.Background())
	}

	ip, err := container.Host(ctx)
	if err != nil {
		tearDown()
		return DBConfig{}, nil, faults.Wrap(err)
	}
	port, err := container.MappedPort(ctx, natPort)
	if err != nil {
		tearDown()
		return DBConfig{}, nil, faults.Wrap(err)
	}

	dbConfig.Host = ip
	dbConfig.Port = port.Int()

	return dbConfig, tearDown, nil
}
