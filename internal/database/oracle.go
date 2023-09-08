package database

import (
	"context"
	"fmt"
	"support/internal/env"

	_ "github.com/godror/godror"
	"github.com/jmoiron/sqlx"
)

// DB Class
type Oracle struct {
	ctx    context.Context
	config env.EnvApp
}

func (o *Oracle) connect() (*sqlx.DB, error) {
	connectionString := fmt.Sprintf(`user="%s" password="%s" connectString="%s"`, o.config.DB_USERNAME, o.config.DB_PASSWORD, fmt.Sprintf("%s:%s/%s", o.config.DB_HOST, o.config.DB_PORT, o.config.DB_SERVICE))
	db, err := sqlx.Open("godror", connectionString)
	if err != nil {
		panic(err)
	}
	var queryResultColumnOne string
	row := db.QueryRow("SELECT systimestamp FROM dual")
	err = row.Scan(&queryResultColumnOne)
	if err != nil {
		panic(fmt.Errorf("error scanning db: %w", err))
	}
	fmt.Println("The time in the database ", queryResultColumnOne)
	return db, nil
}

func (c *Oracle) Connect() *sqlx.DB {
	counts := 0

	for {
		db, err := c.connect()

		if try(err, db, &counts) == nil {
			return db
		}
		continue
	}
}
func NewOracleDB(ctx context.Context, ec env.EnvApp) Database {
	return &Oracle{ctx: ctx, config: ec}
}
