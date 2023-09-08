package cmd

import (
	"context"
	"support/internal/database"
	"support/internal/env"
	"support/internal/server"

	"github.com/ignaciocaff/oraclesp"
)

func Start() {
	ctx := context.Background()
	_env := env.GetEnv(".env")

	db := database.NewOracleDB(ctx, _env).Connect()
	oraclesp.Configure(db, ctx)
	defer db.Close()
	server.NewGinServer(ctx, db, _env).Start()
}
