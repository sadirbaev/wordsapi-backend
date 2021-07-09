package server

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/olivere/elastic/v7"
	"github.com/spf13/viper"
	"time"
	documentDelivery "wordsapi/document/delivery"
	documentRepo "wordsapi/document/repository/esservice"
	documentUsecase "wordsapi/document/usecase"
	"wordsapi/domain"
	"wordsapi/middleware"
)

type App struct {
	documentUsecase domain.DocumentUsecase
}

func NewServer(es *elastic.Client) *App {
	timeoutContext := time.Duration(viper.GetInt("context.timeout")) * time.Second
	newDocumentRepo := documentRepo.NewESDocumentRepository(es)
	return &App{
		documentUsecase: documentUsecase.NewDocumentUsecase(newDocumentRepo, timeoutContext),
	}
}

func (rx *App) Run() error {
	e := echo.New()
	middL := middleware.InitMiddleware()
	e.Use(middL.CORS)
	g := e.Group("/api")
	documentDelivery.NewDocumentHandler(g, rx.documentUsecase)
	fmt.Println(fmt.Sprintf("%s", e.Start(viper.GetString("server.address"))))
	return nil
}
