package main

import (
	"log"
	"redditclone/pkg/handlers"
	"redditclone/pkg/items"
	"redditclone/pkg/server"
	"redditclone/pkg/session"
	"redditclone/pkg/user"

	"go.uber.org/zap"
)

func main() {
	zapLogger, ErrLog := zap.NewProduction()
	if ErrLog != nil {
		log.Fatal("logger error")
	}
	logger := zapLogger.Sugar()

	userHandler := &handlers.UserHandler{
		UserRepo: user.NewUserRepo(),
		Logger:   logger,
		Sessions: session.NewSessionsMem(),
	}

	itemHandler := &handlers.ItemsHandler{
		ItemsRepo: items.NewRepo(),
		Logger:    logger,
	}

	ErrHTTP := server.Run(userHandler, itemHandler, ":8091")
	if ErrHTTP != nil {
		log.Fatal(ErrHTTP)
	}
	ErrLog = zapLogger.Sync()
	if ErrLog != nil {
		log.Fatal("logger error")
	}
}
