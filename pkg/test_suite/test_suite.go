package test_suite

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
	nats2 "github.com/nats-io/nats.go"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/nats"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

type DBTestSuite struct {
	suite.Suite
	DBPool            *pgxpool.Pool
	PostgresContainer testcontainers.Container
	Ctx               context.Context
	Migrations        string
	Connection        string
	Driver            string
	DBName            string

	NatsContainer *nats.NATSContainer
	Nc            *nats2.Conn
}

// NewDBTestSuite creates a new DBTestSuite with default values
func NewDBTestSuite(migrationsPath string) *DBTestSuite {
	return &DBTestSuite{
		Ctx:        context.Background(),
		Migrations: migrationsPath,
		Driver:     "pgx",
		DBName:     "postgres",
	}
}

func (suite *DBTestSuite) SetupSuite() {
	// Use TestContainers instead of embedded-postgres
	ctx := context.Background()

	container, err := postgres.Run(ctx,
		"postgres:15-alpine",
		postgres.WithDatabase("testdb"),
		postgres.WithUsername("admin"),
		postgres.WithPassword("Avanp0st"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2),
		),
	)
	if err != nil {
		panic(fmt.Sprintf("Failed to start postgres postgresContainer: %v", err))
	}

	suite.PostgresContainer = container

	// Get connection string
	connStr, err := container.ConnectionString(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to get connection string: %v", err))
	}

	// Add SSL mode
	suite.Connection = connStr + "&sslmode=disable"

	// Run migrations
	if err := suite.migrateUp(); err != nil {
		panic(fmt.Sprintf("Failed to run migrations: %v", err))
	}

	// Create connection pool
	pool, err := suite.connect()
	if err != nil {
		panic(fmt.Sprintf("Failed to connect to database: %v", err))
	}

	suite.DBPool = pool

	suite.NatsContainer, err = nats.Run(ctx, "nats:2.9")
	uri, err := suite.NatsContainer.ConnectionString(ctx)
	if err != nil {
		panic(fmt.Sprintf("Failed to get nats connection string: %v", err))
	}
	if suite.Nc, err = nats2.Connect(uri); err != nil {
		panic(fmt.Sprintf("Failed to connect nats: %v", err))
	}
}

func (suite *DBTestSuite) TearDownSuite() {
	if suite.DBPool != nil {
		suite.DBPool.Close()
	}

	if suite.PostgresContainer != nil {
		if err := suite.PostgresContainer.Terminate(suite.Ctx); err != nil {
			suite.T().Errorf("Failed to stop postgresContainer: %v", err)
		}
	}

	if suite.Nc != nil {
		_ = suite.Nc.Drain()
		suite.Nc.Close()
	}

	if suite.NatsContainer != nil {
		if err := suite.NatsContainer.Terminate(suite.Ctx); err != nil {
			suite.T().Errorf("Failed to stop postgresContainer: %v", err)
		}
	}
}

func (suite *DBTestSuite) connect() (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(suite.Connection)
	if err != nil {
		return nil, err
	}
	return pgxpool.NewWithConfig(context.TODO(), poolCfg)
}

func (suite *DBTestSuite) migrateUp() error {
	sourceURL := "file://" + suite.Migrations
	conn, err := sql.Open(suite.Driver, suite.Connection)
	if err != nil {
		return err
	}
	defer conn.Close()

	instance, err := pgxmigrate.WithInstance(conn, &pgxmigrate.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(sourceURL, suite.DBName, instance)
	if err != nil {
		return err
	}

	if err = m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			return nil
		}
		return err
	}
	return nil
}

// CleanupTables provides a generic way to clean up tables
func (suite *DBTestSuite) CleanupTables(tables []string) {
	for _, table := range tables {
		_, err := suite.DBPool.Exec(suite.Ctx, "DELETE FROM "+table)
		if err != nil {
			suite.T().Errorf("Failed to clean table %s: %v", table, err)
		}
	}
}
