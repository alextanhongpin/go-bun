package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"github.com/uptrace/bun/extra/bundebug"
)

type PostgresConfig struct {
	Name     string `required:"true"`
	User     string `required:"true"`
	Password string `required:"true" envconfig:"pass"`
	Host     string `required:"true" envconfig:"host"`
	Port     int64  `required:"true"`
}

func (cfg PostgresConfig) String() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
}

func main() {
	var cfg PostgresConfig
	if err := envconfig.Process("db", &cfg); err != nil {
		panic(err)
	}

	// dsn := "unix://user:pass@dbname/var/run/postgresql/.s.PGSQL.5432"
	dsn := cfg.String()
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))

	db := bun.NewDB(sqldb, pgdialect.New())
	defer db.Close()

	// Print queries to stdout.
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	ctx := context.Background()
	res, err := db.ExecContext(ctx, "SELECT 1")
	if err != nil {
		panic(err)
	}
	fmt.Println(res)

	var num int
	if err := db.QueryRowContext(ctx, "SELECT 1").Scan(&num); err != nil {
		panic(err)
	}

	{

		res, err := db.NewSelect().ColumnExpr("1").Exec(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}

	{
		var num int
		if err := db.NewSelect().ColumnExpr("1").Scan(ctx, &num); err != nil {
			panic(err)
		}
		fmt.Println(num)
	}
}
